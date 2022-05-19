package db

import (
	"koudai-box/conf"

	"github.com/astaxie/beego/orm"
	// "github.com/mattn/go-sqlite3" // driver
	"github.com/sirupsen/logrus"
)

var webOrm orm.Ormer
var logger = logrus.New()

func CreateTables() {
	orm.ResetModelCache()
	RegisterModels()
	err := orm.RegisterDataBase("default", "sqlite3", conf.GetConf().DbPath+"/iot.db", 30)
	if err != nil {
		logrus.Errorln(err)
	}
	err = orm.RunSyncdb("default", false, true)
	if err != nil {
		logrus.Errorln(err)
	}
	orm.ResetModelCache()
}

func RegisterModels() {
	orm.RegisterModel(new(Gateway))
	orm.RegisterModel(new(Alarm))
	orm.RegisterModel(new(Event))
	orm.RegisterModel(new(Dict))
	orm.RegisterModel(new(Product))
	orm.RegisterModel(new(DeviceProperty))
	orm.RegisterModel(new(Device))
	orm.RegisterModel(new(GatewayApiHistory))
	orm.RegisterModel(new(OperationRecord))
	orm.RegisterModel(new(DeviceStatus))
}

// Init db
func Init() {
	err := orm.RegisterDataBase("default", "sqlite3", conf.GetConf().DbPath+"/web.db", 30)
	if err != nil {
		logrus.Errorln(err)
	}
	webOrm = orm.NewOrm()
	webOrm.Using("default")
}

func InitDict() {
	logrus.Println("init dict db")
	sqls := []string{
		`INSERT INTO "main"."dict"("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (1, 0, 'gateway_protocol_type', 0, 'HTTP客户端', 'http_client', '', '');`,
		`INSERT INTO "main"."dict"("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (2, 0, 'gateway_protocol_type', 1, 'HTTP服务端', 'http_server', '', '');`,
		`INSERT INTO "main"."dict"("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (3, 0, 'gateway_status', 2, '运行中', 'running', '', '');`,
		`INSERT INTO "main"."dict"("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (4, 0, 'gateway_status', 3, '停止', 'stop', '', '');`,
		`INSERT INTO "main"."dict"("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (5, 0, 'http_request_type', 4, 'GET', 'get', '', '');`,
		`INSERT INTO "main"."dict"("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (6, 0, 'http_request_type', 5, 'POST', 'post', '', '');`,
		`INSERT INTO "main"."dict"("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (7, 0, 'data_combination', 6, '单条', 'single', '', '');`,
		`INSERT INTO "main"."dict"("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (8, 0, 'data_combination', 7, '批量', 'array', '', '');`,
		`INSERT INTO "main"."dict"("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (9, 0, 'collect_type', 8, '定时', 'schedule', '', '');`,
		`INSERT INTO "main"."dict"("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (10, 0, 'collect_type', 9, '轮询', 'poll', '', '');`,

		`INSERT INTO "main"."dict" ("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (11, 0, 'product_category', 10, '传感器', 'sensor', '', '');`,
		`INSERT INTO "main"."dict" ("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (12, 11, 'product_category', 11, '普通传感器', 'common_sensor1', '', '');`,
		`INSERT INTO "main"."dict" ("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (13, 11, 'product_category', 12, '智能传感器', 'intelligent_sensor2', '', '');`,

		`INSERT INTO "main"."dict"("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (15, 0, 'field_type', 14, '整型', '1', '', '');`,
		`INSERT INTO "main"."dict"("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (16, 0, 'field_type', 15, '枚举', '4', '', '');`,
		`INSERT INTO "main"."dict"("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (17, 0, 'field_type', 16, '结构体', '9', '', '');`,
		`INSERT INTO "main"."dict"("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (18, 0, 'field_type', 17, '长整数型', '13', '', '');`,
		`INSERT INTO "main"."dict"("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (19, 0, 'field_type', 18, '单精度浮点数', '2', '', '');`,
		`INSERT INTO "main"."dict"("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (20, 0, 'field_type', 19, '双精度浮点数', '3', '', '');`,
		`INSERT INTO "main"."dict"("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (21, 0, 'field_type', 20, '字符串', '6', '', '');`,
		`INSERT INTO "main"."dict"("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (22, 0, 'field_type', 21, '布尔型', '5', '', '');`,
		`INSERT INTO "main"."dict"("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (23, 0, 'field_type', 22, '时间型', '7', '', '');`,
		`INSERT INTO "main"."dict"("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (24, 0, 'field_type', 23, '数组', '10', '', '');`,
		`INSERT INTO "main"."dict"("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (25, 0, 'field_type', 24, '文件', '11', '', '');`,
		`INSERT INTO "main"."dict"("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (26, 0, 'field_type', 25, '密码', '12', '', '');`,
		`INSERT INTO "main"."dict"("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (27, 0, 'field_type', 26, '地理位置', '14', '', '');`,
		`INSERT INTO "main"."dict" ("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (28, 0, 'modbus_operations', 27, '读线圈状态（0x01）', '0x01', '', '');`,
		`INSERT INTO "main"."dict" ("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (29, 0, 'modbus_operations', 28, '读输入离散量（0x02）', '0x02', '', '');`,
		`INSERT INTO "main"."dict" ("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (30, 0, 'modbus_operations', 29, '读多个寄存器（0x03）', '0x03', '', '');`,
		`INSERT INTO "main"."dict" ("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (31, 0, 'modbus_operations', 30, '读输入寄存器（0x04）', '0x04', '', '');`,
		`INSERT INTO "main"."dict" ("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (32, 0, 'modbus_operations', 31, '写单个线圈（0x05）', '0x05', '', '');`,
		`INSERT INTO "main"."dict" ("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (33, 0, 'modbus_operations', 32, '写单个保持寄存器（0x06）', '0x06', '', '');`,
		`INSERT INTO "main"."dict" ("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (34, 0, 'modbus_operations', 33, '读取异常状态（0x07）', '0x07', '', '');`,
		`INSERT INTO "main"."dict" ("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (35, 0, 'modbus_operations', 34, '回送诊断校验（0x08）', '0x08', '', '');`,
		`INSERT INTO "main"."dict" ("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (36, 0, 'modbus_operations', 35, '编程（只用于484）(0x09)', '0x09', '', '');`,
		`INSERT INTO "main"."dict" ("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (37, 0, 'modbus_operations', 36, '控询（只用于484）(0x0A)', '0x0A', '', '');`,
		`INSERT INTO "main"."dict" ("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (38, 0, 'modbus_operations', 37, '读取事件计数（0x0B）', '0x0B', '', '');`,
		`INSERT INTO "main"."dict" ("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (39, 0, 'modbus_operations', 38, '读取通讯事件记录（0x0C）', '0x0C', '', '');`,
		`INSERT INTO "main"."dict" ("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (40, 0, 'modbus_operations', 39, '编程（184/384/484/584）(0x0D)', '0x0D', '', '');`,
		`INSERT INTO "main"."dict" ("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (41, 0, 'modbus_operations', 40, '探询（184/384/484/584）(0x0E)', '0x0E', '', '');`,
		`INSERT INTO "main"."dict" ("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (42, 0, 'modbus_operations', 41, '写多个线圈（0x0F）', '0x0F', '', '');`,
		`INSERT INTO "main"."dict" ("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (43, 0, 'modbus_operations', 42, '写多个保持寄存器（0x10)', '0x10', '', '');`,
		`INSERT INTO "main"."dict" ("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (44, 0, 'modbus_operations', 43, '报告从机标识（0x011）', '0x11', '', '');`,
		`INSERT INTO "main"."dict" ("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (45, 0, 'modbus_operations', 44, '（884和MICRO84）(0x12)', '0x12', '', '');`,
		`INSERT INTO "main"."dict" ("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (46, 0, 'modbus_operations', 45, '重置通信链路(0x13)', '0x13', '', '');`,
		`INSERT INTO "main"."dict" ("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (47, 0, 'modbus_operations', 46, '读取通用参数（584L）(0x14)', '0x14', '', '');`,
		`INSERT INTO "main"."dict" ("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (48, 0, 'modbus_operations', 47, '写入通用参数（584L）(0x15)', '0x15', '', '');`,
		`INSERT INTO "main"."dict"("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (49, 0, 'gateway_protocol_type', 48, 'Modbus RTU', 'Modbus_RTU', '', '');`,
		`INSERT INTO "main"."dict"("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (50, 0, 'gateway_protocol_type', 49, 'Modbus TCP', 'Modbus_TCP', '', '');`,
		`INSERT INTO "main"."dict"("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (51, 0, 'gateway_protocol_type', 50, 'OPC UA', 'OPC_UA', '', '');`,
		`INSERT INTO "main"."dict"("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (52, 0, 'gateway_protocol_type', 51, 'BACnet/IP', 'BACnet_IP', '', '');`,
		`INSERT INTO "main"."dict"("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (53, 0, 'gateway_protocol_type', 52, 'MQTT', 'MQTT', '', '');`,
		`INSERT INTO "main"."dict"("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (54, 0, 'gateway_protocol_type', 53, 'MQTT-SN', 'MQTT-SN', '', '');`,
		`INSERT INTO "main"."dict"("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (55, 0, 'gateway_protocol_type', 54, 'CoAP', 'CoAP', '', '');`,
		`INSERT INTO "main"."dict"("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (56, 0, 'gateway_protocol_type', 55, 'LwM2M', 'LwM2M', '', '');`,
		`INSERT INTO "main"."dict"("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (57, 0, 'gateway_protocol_type', 56, 'WebSocket服务端', 'Websocket_server', '', '');`,
		`INSERT INTO "main"."dict"("id", "pid", "type", "sort", "name", "value", "tips", "extra") VALUES (58, 0, 'gateway_protocol_type', 57, 'WebSocket客户端', 'Websocket_client', '', '');`,
	}
	for _, s := range sqls {
		ss := s
		_, err := webOrm.Raw(ss).Exec()
		if err != nil {
			continue
		}
	}
}
