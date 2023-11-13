test_http:
	go test -count=1 -v github.com/smiecj/go_common/http -run="TestMakeGetRequest"

test_db_memory:
	go test -count=1 -v github.com/smiecj/go_common/db/local -run="TestLocalMemoryConnector"

test_db_file:
	go test -count=1 -v github.com/smiecj/go_common/db/local -run="TestLocalFileConnector"

test_db_mysql:
	go test -count=1 -v github.com/smiecj/go_common/db/mysql -run="TestMySQLConnector"

test_db_batch:
	go test -count=1 -v github.com/smiecj/go_common/db/mysql -run="TestMySQLBatchInsert"

test_db_impala:
	go test -count=1 -v github.com/smiecj/go_common/db/impala -run="TestImpalaConnector"

test_zk:
	go test -count=1 -v github.com/smiecj/go_common/zk -run="TestZKConnect"

test_yaml_config:
	go test -count=1 -v github.com/smiecj/go_common/config/yaml -run="TestYamlConfig"

test_nacos_config:
	go test -count=1 -v github.com/smiecj/go_common/config/nacos -run="TestNacosConfig"

test_apollo_config:
	go test -count=1 -timeout=100s -v github.com/smiecj/go_common/config/apollo -run="TestApolloConfig"

test_smtp:
	go test -timeout=10s -count=1 -v github.com/smiecj/go_common/util/mail -run="TestSendMail"