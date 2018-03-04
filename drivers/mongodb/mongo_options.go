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
		wantCommandString string
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
			wantCommandString: "mongodump --host 127.0.0.1 --port 27017 --out=/opt/dump -u test -p test -j=10 --ipv6 --gzip --db=testDB",
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
			wantCommandString: "mongodump --host 127.0.0.1 --port 27017 --out=/opt/dump -u test -p test -j=10 --gzip --db=testDB",
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
			wantCommandString: "mongodump --host 127.0.0.1 --port 27017 --out=/opt/dump -u test -p test -j=10 --ipv6 --db=testDB",
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
			wantCommandString: "mongodump --host 127.0.0.1 --port 27017 --out=/opt/dump -u test -p test -j=10 --db=testDB",
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
			wantCommandString: "mongodump --port 27017 --out=/opt/dump -u test -p test -j=10 --ipv6 --gzip --db=testDB",
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
			wantCommandString: "mongodump --host 127.0.0.1 --out=/opt/dump -u test -p test -j=10 --ipv6 --gzip --db=testDB",
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
			wantCommandString: "mongodump --host 127.0.0.1 --port 27017 --out=/opt/dump -j=10 --ipv6 --gzip --db=testDB",
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
			wantCommandString: "mongodump --host 127.0.0.1 --port 27017 --out=/opt/dump -u test -p test -j=10 --ipv6 --gzip",
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
			wantCommandString: "mongodump --host 127.0.0.1 --port 27017 --out=/opt/dump -u test -p test --ipv6 --gzip --db=testDB",
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
			wantCommandString: "mongodump --host 127.0.0.1 --port 27017 -u test -p test -j=10 --ipv6 --gzip --db=testDB",
		},
		{
			name: "Minumum",
			fields: fields{
				cfg: options.InitControlConfig(),
			},
			wantCommandString: "mongodump",
		},
		{
			name: "With Host Port and Database",
			fields: fields{
				cfg:      options.InitControlConfig(),
				host:     "127.0.0.1",
				port:     27017,
				database: "testDB",
			},
			wantCommandString: "mongodump --host 127.0.0.1 --port 27017 --db=testDB",
		},
	}
)
