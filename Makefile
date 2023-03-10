test_http:
	go test -count=1 -v github.com/smiecj/go_common/http -run="TestMakeGetRequest"

test_db_memory:
	go test -count=1 -v github.com/smiecj/go_common/db/local -run="TestLocalMemoryConnector"

test_db_file:
	go test -count=1 -v github.com/smiecj/go_common/db/local -run="TestLocalFileConnector"

test_db_mysql:
	go test -count=1 -v github.com/smiecj/go_common/db/mysql -run="TestMySQLConnector" -config=/tmp/conf.yaml

test_db_impala:
	go test -count=1 -v github.com/smiecj/go_common/db/impala -run="TestImpalaConnector"

test_config:
	go test -count=1 -v github.com/smiecj/go_common/config -run="TestYamlConfig"

test_smtp:
	go test -count=1 -v github.com/smiecj/go_common/util/mail -run="TestSendMail"