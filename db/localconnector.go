package db

type localMemoryConnector struct {
	// db_name -> key: table name; id: uuid
	storage map[string]map[string]string
}

func (connector *localMemoryConnector) init() {
	connector.storage = make(map[string]map[string]string)
}

// 本地存储: 插入数据
func (connector *localMemoryConnector) Insert(funcArr ...rdbInsertConfigFunc) (updateRet, error) {
	action := makeRDBAddAction()
	for _, currentFunc := range funcArr {
		currentFunc(action)
	}

	// 本地存储仅会存入 key-value 格式的数据，并且会覆盖
	spaceName := action.getSpaceName()
	if nil == connector.storage[spaceName] {
		connector.storage[spaceName] = make(map[string]string)
	}
	for _, currentField := range action.fieldArr {
		connector.storage[spaceName][currentField.key] = currentField.value
	}

	return updateRet{AffectedRows: len(action.fieldArr)}, nil
}

// todo: 实现本地内存连接器
func GetLocalMemoryConnector() *RDBConnector {
	return nil
}

// 后续: 初始化 连接器配置中，增加 id generator 配置
