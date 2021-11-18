test_http:
	go test -count=1 -v github.com/smiecj/go_common/http -run="TestMakeGetRequest"

test_db_memory:
	go test -count=1 -v github.com/smiecj/go_common/db -run="TestLocalMemoryConnector"

test_db_file:
	go test -count=1 -v github.com/smiecj/go_common/db -run="TestLocalFileConnector"

test_db_mysql:
	go test -count=1 -v github.com/smiecj/go_common/db -run="TestMySQLConnector"

test_config:
	go test -count=1 -v github.com/smiecj/go_common/config -run="TestYamlConfig"