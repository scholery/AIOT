package db

import (
	"github.com/astaxie/beego/orm"
)

//添加数据
func InsertProduct(product Product) (int64, error) {
	id, err := webOrm.Insert(&product)
	if err != nil {
		logger.Errorln(err)
		id = -1
	}
	return id, err
}

func QueryProductByIDs(productIds []int) []*Product {
	var products []*Product
	qs := webOrm.QueryTable("product").Filter("del_flag", 0).Filter("id__in", productIds)

	qs.All(&products)
	return products
}

//查询数据库
func QueryProductByPage(offset, limit int, search string, state int) (int64, []*Product) {
	var products []*Product
	qs := webOrm.QueryTable("product")
	cond := orm.NewCondition()
	cond = cond.And("del_flag", 0)

	if state != -1 {
		cond = cond.And("state", state)
	}

	if len(search) > 0 {
		cond1 := orm.NewCondition()
		cond1 = cond1.Or("code__contains", search).Or("name__contains", search)
		cond = cond.AndCond(cond1)
	}

	qs = qs.SetCond(cond)
	_, err := qs.Limit(limit, offset).OrderBy("-createTime").All(&products)
	if err != nil {
		logger.Errorln(err)
	}
	num, _ := qs.Count()
	return num, products
}

//查询所有产品数量
func QueryAllProductCount() int64 {
	num, err := webOrm.QueryTable("product").Filter("del_flag", 0).Count()
	if err != nil {
		logger.Errorln(err)
		return 0
	}
	return num
}

//查询所有产品数量
func QueryAllProductByStateCount(state int) int64 {
	num, err := webOrm.QueryTable("product").Filter("state", state).Filter("del_flag", 0).Count()
	if err != nil {
		logger.Errorln(err)
		return 0
	}
	return num
}

//更加id查询
func QueryProductByID(id int) (*Product, error) {
	var product Product
	err := webOrm.QueryTable("product").Filter("id", id).Filter("del_flag", 0).One(&product)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	return &product, nil
}

//根据网关id查询
func QueryProductByGatewayID(gatewayId int) ([]*Product, error) {
	var products []*Product
	_, err := webOrm.QueryTable("product").Filter("gatewayId", gatewayId).Filter("del_flag", 0).All(&products)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	return products, nil
}

//更新
func UpdateProduct(product *Product) error {
	webOrm.Begin()
	_, err := webOrm.Update(product)
	if err != nil {
		logger.Errorln(err)
		webOrm.Rollback()
	}
	webOrm.Commit()
	return err
}
