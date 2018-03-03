package mongodb

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"shkaff/internal/consts"
	"shkaff/internal/databases"
	"shkaff/internal/options"
	"shkaff/internal/structs"
	"strings"
)

var (
	MONGO_SUCESS_DUMP = regexp.MustCompile(`\tdone\ dumping`)
	RESTORE_ERRORS    = [2]*regexp.Regexp{
		regexp.MustCompile(`exit\ status\ 1`),
		regexp.MustCompile(`skipping...`)}
)

type MongoParams struct {
	cfg                    *options.ShkaffConfig
	host                   string
	port                   int
	login                  string
	password               string
	ipv6                   bool
	database               string
	collection             string
	gzip                   bool
	parallelCollectionsNum int
	dumpFolder             string
	resultChan             chan string
}

func InitDriver() (mp databases.DatabaseDriver) {
	return &MongoParams{
		cfg: options.InitControlConfig(),
	}
}

func (mp *MongoParams) setDBSettings(task *structs.Task) {
	mp.host = task.Host
	mp.port = task.Port
	mp.login = task.DBUser
	mp.password = task.DBPassword
	mp.ipv6 = task.Ipv6
	mp.gzip = task.Gzip
	mp.database = task.Database
	mp.parallelCollectionsNum = task.ThreadCount
	mp.dumpFolder = task.DumpFolder

}

func (mp *MongoParams) isUseAuth() bool {
	return mp.login != "" && mp.password != ""
}

func (mp *MongoParams) ParamsToDumpString() (commandString string) {
	var cmdLine []string

	cmdLine = append(cmdLine, consts.MONGO_DUMP_COMMAND)
	cmdLine = append(cmdLine, fmt.Sprintf("%s %s", consts.MONGO_HOST_KEY, mp.host))
	cmdLine = append(cmdLine, fmt.Sprintf("%s %d", consts.MONGO_PORT_KEY, mp.port))
	cmdLine = append(cmdLine, fmt.Sprintf("--out=%s", mp.dumpFolder))

	if mp.isUseAuth() {
		auth := fmt.Sprintf("%s %s %s %s", consts.MONGO_LOGIN_KEY, mp.login, consts.MONGO_PASS_KEY, mp.password)
		cmdLine = append(cmdLine, auth)
	}

	if mp.ipv6 {
		cmdLine = append(cmdLine, consts.MONGO_GZIP_KEY)
	}

	if mp.gzip {
		cmdLine = append(cmdLine, consts.MONGO_GZIP_KEY)
	}

	if mp.database != "" {
		cmdLine = append(cmdLine, fmt.Sprintf("%s=%s", consts.MONGO_DATABASE_KEY, mp.database))
		if mp.collection != "" {
			cmdLine = append(cmdLine, fmt.Sprintf("%s=%s", consts.MONGO_COLLECTION_KEY, mp.collection))
		}
	}
	if mp.collection == "" && mp.parallelCollectionsNum > 4 {
		cmdLine = append(cmdLine, fmt.Sprintf("%s=%d", consts.MONGO_PARALLEL_KEY, mp.parallelCollectionsNum))
	}

	commandString = strings.Join(cmdLine, " ")
	return
}

func (mp *MongoParams) Dump(task *structs.Task) (err error) {
	var stderr bytes.Buffer
	mp.setDBSettings(task)
	cmd := exec.Command("sh", "-c", mp.ParamsToDumpString())
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	dumpResult := stderr.String()
	reResult := MONGO_SUCESS_DUMP.FindString(dumpResult)
	if reResult != "" {
		return
	}
	return errors.New("Restore: " + dumpResult)
}
func (mp *MongoParams) ParamsToRestoreString() (commandString string) {
	var cmdLine []string
	cmdLine = append(cmdLine, "mongorestore")
	cmdLine = append(cmdLine, fmt.Sprintf("%s %s", consts.MONGO_HOST_KEY, mp.cfg.MONGO_RESTORE_HOST))
	cmdLine = append(cmdLine, fmt.Sprintf("%s %d", consts.MONGO_PORT_KEY, mp.cfg.MONGO_RESTORE_PORT))
	if mp.isUseAuth() {
		auth := fmt.Sprintf("%s %s %s %s", consts.MONGO_LOGIN_KEY, mp.login, consts.MONGO_PASS_KEY, mp.password)
		cmdLine = append(cmdLine, auth)
	}
	if mp.ipv6 {
		cmdLine = append(cmdLine, consts.MONGO_GZIP_KEY)
	}
	if mp.gzip {
		cmdLine = append(cmdLine, consts.MONGO_GZIP_KEY)
	}
	//	if mp.collection == "" && mp.parallelCollectionsNum > 4 {
	//		cmdLine = append(cmdLine, fmt.Sprintf("%s=%d", consts.MONGO_PARALLEL_KEY, mp.parallelCollectionsNum))
	//	}
	dir := fmt.Sprintf("-d %s '%s/%s'", mp.database, mp.dumpFolder, mp.database)
	cmdLine = append(cmdLine, dir)
	cmdLine = append(cmdLine, "--stopOnError")
	cmdLine = append(cmdLine, "--drop")
	commandString = strings.Join(cmdLine, " ")
	return
}

func (mp *MongoParams) Restore(task *structs.Task) (err error) {
	var stderr bytes.Buffer
	mp.setDBSettings(task)
	cmd := exec.Command("sh", "-c", mp.ParamsToRestoreString())
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	restoreResult := stderr.String()
	for _, restoreErrorPattern := range RESTORE_ERRORS {
		reResult := restoreErrorPattern.FindString(restoreResult)
		if reResult != "" {
			return errors.New("Restore: " + strings.TrimSpace(restoreResult))
		}
	}
	return
}
