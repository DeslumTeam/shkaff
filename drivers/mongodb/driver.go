package mongodb

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/DeslumTeam/shkaff/internal/consts"
	"github.com/DeslumTeam/shkaff/internal/databases"
	"github.com/DeslumTeam/shkaff/internal/options"
	"github.com/DeslumTeam/shkaff/internal/structs"
)

type MongoParams struct {
	cfg                    *options.ShkaffConfig
	host                   string
	port                   int
	login                  string
	password               string
	ipv6                   bool
	database               string
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
	if mp.host != "" {
		cmdLine = append(cmdLine, fmt.Sprintf("%s %s", consts.MONGO_HOST_KEY, mp.host))
	}

	if mp.port > 0 {
		cmdLine = append(cmdLine, fmt.Sprintf("%s %d", consts.MONGO_PORT_KEY, mp.port))
	}
	if mp.dumpFolder != "" {
		cmdLine = append(cmdLine, fmt.Sprintf("--out=%s", mp.dumpFolder))
	}

	if mp.isUseAuth() {
		// TODO admin is different
		auth := fmt.Sprintf("%s %s %s %s %s=admin", consts.MONGO_LOGIN_KEY,
			mp.login,
			consts.MONGO_PASS_KEY,
			mp.password,
			consts.MONGO_AUTH_DB_KEY)
		cmdLine = append(cmdLine, auth)
	}

	if mp.parallelCollectionsNum > 4 {
		cmdLine = append(cmdLine, fmt.Sprintf("%s=%d", mp.parallelCollectionsNum))
	}

	if mp.ipv6 {
		cmdLine = append(cmdLine, consts.MONGO_IPV6_KEY)
	}

	if mp.gzip {
		cmdLine = append(cmdLine, consts.MONGO_GZIP_KEY)
	}

	if mp.database != "" {
		cmdLine = append(cmdLine, fmt.Sprintf("%s=%s", consts.MONGO_DATABASE_KEY, mp.database))
	}
	cmdLine = append(cmdLine, "-v")
	commandString = strings.Join(cmdLine, " ")
	return
}

func (mp *MongoParams) Dump(task *structs.Task) (err error) {
	mp.setDBSettings(task)
	log.Println(mp.ParamsToDumpString())
	cmd := exec.Command("sh", "-c", mp.ParamsToDumpString())
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		dumpResult := scanner.Text()
		log.Println(dumpResult)
	}
	err = cmd.Wait()
	return
}
func (mp *MongoParams) ParamsToRestoreString() (commandString string) {
	var cmdLine []string
	cmdLine = append(cmdLine, "mongorestore")
	if mp.host != "" {
		cmdLine = append(cmdLine, fmt.Sprintf("%s %s", consts.MONGO_HOST_KEY, mp.cfg.MONGO_RESTORE_HOST))
	}
	if mp.port > 0 {
		cmdLine = append(cmdLine, fmt.Sprintf("%s %d", consts.MONGO_PORT_KEY, mp.cfg.MONGO_RESTORE_PORT))
	}

	if mp.isUseAuth() {
		// TODO admin is different
		auth := fmt.Sprintf("%s %s %s %s %s=admin", consts.MONGO_LOGIN_KEY,
			mp.login,
			consts.MONGO_PASS_KEY,
			mp.password,
			consts.MONGO_AUTH_DB_KEY)
		cmdLine = append(cmdLine, auth)
	}

	if mp.ipv6 {
		cmdLine = append(cmdLine, consts.MONGO_IPV6_KEY)
	}
	if mp.gzip {
		cmdLine = append(cmdLine, consts.MONGO_GZIP_KEY)
	}
	if mp.database != "" {
		db := fmt.Sprintf("--nsInclude=%s.*", mp.database)
		cmdLine = append(cmdLine, db)
	}
	if mp.dumpFolder != "" {
		restorePath := fmt.Sprintf("--dir=%s", mp.dumpFolder)
		cmdLine = append(cmdLine, restorePath)
	}

	cmdLine = append(cmdLine, "--drop -v")
	commandString = strings.Join(cmdLine, " ")
	return
}

func (mp *MongoParams) Restore(task *structs.Task) (err error) {
	mp.setDBSettings(task)
	log.Println(mp.ParamsToRestoreString())
	cmd := exec.Command("sh", "-c", mp.ParamsToRestoreString())
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		restoreResult := scanner.Text()
		log.Println(restoreResult)
	}
	err = cmd.Wait()
	return
}
