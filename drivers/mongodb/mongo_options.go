package mongodb

import "github.com/DeslumTeam/shkaff/internal/options"

type fields struct {
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

var (
	tests = []struct {
		name              string
		fields            fields
		wantDumpString    string
		wantRestoreString string
	}{
		{
			name: "Full",
			fields: fields{
				cfg:        options.InitControlConfig(),
				host:       "127.0.0.1",
				port:       27017,
				login:      "test",
				password:   "test",
				ipv6:       true,
				database:   "testDB",
				collection: "",
				gzip:       true,
				parallelCollectionsNum: 10,
				dumpFolder:             "/opt/dump",
			},
			wantDumpString:    "mongodump --host 127.0.0.1 --port 27017 --out=/opt/dump -u test -p test --authenticationDatabase=admin -j=10 --ipv6 --gzip --db=testDB -v",
			wantRestoreString: "mongorestore --host mongodb-restore --port 27018 -u test -p test --authenticationDatabase=admin --ipv6 --gzip --nsInclude=testDB.* --dir=/opt/dump --drop -v",
		},
		{
			name: "ipv6 disable",
			fields: fields{
				cfg:        options.InitControlConfig(),
				host:       "127.0.0.1",
				port:       27017,
				login:      "test",
				password:   "test",
				ipv6:       false,
				database:   "testDB",
				collection: "",
				gzip:       true,
				parallelCollectionsNum: 10,
				dumpFolder:             "/opt/dump",
			},
			wantDumpString:    "mongodump --host 127.0.0.1 --port 27017 --out=/opt/dump -u test -p test --authenticationDatabase=admin -j=10 --gzip --db=testDB -v",
			wantRestoreString: "mongorestore --host mongodb-restore --port 27018 -u test -p test --authenticationDatabase=admin --gzip --nsInclude=testDB.* --dir=/opt/dump --drop -v",
		},
		{
			name: "gzip disable",
			fields: fields{
				cfg:        options.InitControlConfig(),
				host:       "127.0.0.1",
				port:       27017,
				login:      "test",
				password:   "test",
				ipv6:       true,
				database:   "testDB",
				collection: "",
				gzip:       false,
				parallelCollectionsNum: 10,
				dumpFolder:             "/opt/dump",
			},
			wantDumpString:    "mongodump --host 127.0.0.1 --port 27017 --out=/opt/dump -u test -p test --authenticationDatabase=admin -j=10 --ipv6 --db=testDB -v",
			wantRestoreString: "mongorestore --host mongodb-restore --port 27018 -u test -p test --authenticationDatabase=admin --ipv6 --nsInclude=testDB.* --dir=/opt/dump --drop -v",
		},
		{
			name: "gzip and ipv6 disable",
			fields: fields{
				cfg:        options.InitControlConfig(),
				host:       "127.0.0.1",
				port:       27017,
				login:      "test",
				password:   "test",
				ipv6:       false,
				database:   "testDB",
				collection: "",
				gzip:       false,
				parallelCollectionsNum: 10,
				dumpFolder:             "/opt/dump",
			},
			wantDumpString:    "mongodump --host 127.0.0.1 --port 27017 --out=/opt/dump -u test -p test --authenticationDatabase=admin -j=10 --db=testDB -v",
			wantRestoreString: "mongorestore --host mongodb-restore --port 27018 -u test -p test --authenticationDatabase=admin --nsInclude=testDB.* --dir=/opt/dump --drop -v",
		},
		{
			name: "Without Host",
			fields: fields{
				cfg:        options.InitControlConfig(),
				port:       27017,
				login:      "test",
				password:   "test",
				ipv6:       true,
				database:   "testDB",
				collection: "",
				gzip:       true,
				parallelCollectionsNum: 10,
				dumpFolder:             "/opt/dump",
			},
			wantDumpString:    "mongodump --port 27017 --out=/opt/dump -u test -p test --authenticationDatabase=admin -j=10 --ipv6 --gzip --db=testDB -v",
			wantRestoreString: "mongorestore --port 27018 -u test -p test --authenticationDatabase=admin --ipv6 --gzip --nsInclude=testDB.* --dir=/opt/dump --drop -v",
		},
		{
			name: "Without port",
			fields: fields{
				cfg:        options.InitControlConfig(),
				host:       "127.0.0.1",
				login:      "test",
				password:   "test",
				ipv6:       true,
				database:   "testDB",
				collection: "",
				gzip:       true,
				parallelCollectionsNum: 10,
				dumpFolder:             "/opt/dump",
			},
			wantDumpString:    "mongodump --host 127.0.0.1 --out=/opt/dump -u test -p test --authenticationDatabase=admin -j=10 --ipv6 --gzip --db=testDB -v",
			wantRestoreString: "mongorestore --host mongodb-restore -u test -p test --authenticationDatabase=admin --ipv6 --gzip --nsInclude=testDB.* --dir=/opt/dump --drop -v",
		},
		{
			name: "Without host & port",
			fields: fields{
				cfg:        options.InitControlConfig(),
				login:      "test",
				password:   "test",
				ipv6:       true,
				database:   "testDB",
				collection: "",
				gzip:       true,
				parallelCollectionsNum: 10,
				dumpFolder:             "/opt/dump",
			},
			wantDumpString:    "mongodump --out=/opt/dump -u test -p test --authenticationDatabase=admin -j=10 --ipv6 --gzip --db=testDB -v",
			wantRestoreString: "mongorestore -u test -p test --authenticationDatabase=admin --ipv6 --gzip --nsInclude=testDB.* --dir=/opt/dump --drop -v",
		},
		{
			name: "Without login password",
			fields: fields{
				cfg:        options.InitControlConfig(),
				host:       "127.0.0.1",
				port:       27017,
				ipv6:       true,
				database:   "testDB",
				collection: "",
				gzip:       true,
				parallelCollectionsNum: 10,
				dumpFolder:             "/opt/dump",
			},
			wantDumpString:    "mongodump --host 127.0.0.1 --port 27017 --out=/opt/dump -j=10 --ipv6 --gzip --db=testDB -v",
			wantRestoreString: "mongorestore --host mongodb-restore --port 27018 --ipv6 --gzip --nsInclude=testDB.* --dir=/opt/dump --drop -v",
		},
		{
			name: "without database",
			fields: fields{
				cfg:        options.InitControlConfig(),
				host:       "127.0.0.1",
				port:       27017,
				login:      "test",
				password:   "test",
				ipv6:       true,
				collection: "",
				gzip:       true,
				parallelCollectionsNum: 10,
				dumpFolder:             "/opt/dump",
			},
			wantDumpString:    "mongodump --host 127.0.0.1 --port 27017 --out=/opt/dump -u test -p test --authenticationDatabase=admin -j=10 --ipv6 --gzip -v",
			wantRestoreString: "mongorestore --host mongodb-restore --port 27018 -u test -p test --authenticationDatabase=admin --ipv6 --gzip --dir=/opt/dump --drop -v",
		},
		{
			name: "Without parallelCollectionsNum",
			fields: fields{
				cfg:                    options.InitControlConfig(),
				host:                   "127.0.0.1",
				port:                   27017,
				login:                  "test",
				password:               "test",
				ipv6:                   true,
				database:               "testDB",
				collection:             "",
				gzip:                   true,
				dumpFolder:             "/opt/dump",
				parallelCollectionsNum: 3,
			},
			wantDumpString:    "mongodump --host 127.0.0.1 --port 27017 --out=/opt/dump -u test -p test --authenticationDatabase=admin --ipv6 --gzip --db=testDB -v",
			wantRestoreString: "mongorestore --host mongodb-restore --port 27018 -u test -p test --authenticationDatabase=admin --ipv6 --gzip --nsInclude=testDB.* --dir=/opt/dump --drop -v",
		},
		{
			name: "Without DumpFolder",
			fields: fields{
				cfg:        options.InitControlConfig(),
				host:       "127.0.0.1",
				port:       27017,
				login:      "test",
				password:   "test",
				ipv6:       true,
				database:   "testDB",
				collection: "",
				gzip:       true,
				parallelCollectionsNum: 10,
			},
			wantDumpString:    "mongodump --host 127.0.0.1 --port 27017 -u test -p test --authenticationDatabase=admin -j=10 --ipv6 --gzip --db=testDB -v",
			wantRestoreString: "mongorestore --host mongodb-restore --port 27018 -u test -p test --authenticationDatabase=admin --ipv6 --gzip --nsInclude=testDB.* --drop -v",
		},
		{
			name: "Minumum",
			fields: fields{
				cfg: options.InitControlConfig(),
			},
			wantDumpString:    "mongodump -v",
			wantRestoreString: "mongorestore --drop -v",
		},
		{
			name: "With Host Port and Database",
			fields: fields{
				cfg:      options.InitControlConfig(),
				host:     "127.0.0.1",
				port:     27017,
				database: "testDB",
			},
			wantDumpString:    "mongodump --host 127.0.0.1 --port 27017 --db=testDB -v",
			wantRestoreString: "mongorestore --host mongodb-restore --port 27018 --nsInclude=testDB.* --drop -v",
		},
	}
)
