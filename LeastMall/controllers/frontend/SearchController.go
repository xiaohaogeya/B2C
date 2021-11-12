package frontend

import (
	"LeastMall/models"
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"github.com/prometheus/common/log"
	"math"
	"reflect"
	"strconv"
)

type SearchController struct {
	BaseController
}

// AddProduct 增加商品数据
func (c *SearchController) AddProduct() {
	var product []models.Product
	models.DB.Find(&product)

	for i := 0; i < len(product); i++ {
		_, err := models.EsClient.Index().
			Index("product").
			Id(strconv.Itoa(product[i].Id)).
			BodyJson(product[i]).
			Do(context.Background())
		if err != nil {
			log.Error(err)
		}
	}

	c.Ctx.WriteString("AddProduct success")
}

// Update 更新数据
func (c *SearchController) Update() {
	//从数据库获取修改
	product := models.Product{}
	models.DB.Where("id=20").Find(&product)
	product.Title = "苹果电脑"
	product.SubTitle = "苹果电脑"
	res, err := models.EsClient.Update().
		Index("product").
		Id("20").
		Doc(product).
		Do(context.Background())
	if err != nil {
		log.Error(err)
	}
	log.Info("update -->", res.Result)
	c.Ctx.WriteString("修改数据")
}

// Delete 删除
func (c *SearchController) Delete() {
	res, err := models.EsClient.Delete().
		Index("product").
		Id("20").
		Do(context.Background())

	if err != nil {
		log.Error(err)
	}
	log.Info("Delete-->", res.Result)

	c.Ctx.WriteString("删除成功")
}

// GetOne 查询一条数据
func (c *SearchController) GetOne() {
	defer func() {
		if r := recover(); r != nil {
			log.Info("Recovered in f", r)
			c.Ctx.WriteString("GetOne")
		}
	}()

	result, _ := models.EsClient.Get().
		Index("product").
		Id("19").
		Do(context.Background())

	log.Info(result.Source)

	product := models.Product{}
	_ = json.Unmarshal(result.Source, &product)
	c.Data["json"] = product
	_ = c.ServeJSON()
}

// Query 查询多条数据
func (c *SearchController) Query() {

	defer func() {
		if r := recover(); r != nil {
			log.Info("Recovered in f", r)
			c.Ctx.WriteString("Query")
		}
	}()

	query := elastic.NewMatchQuery("Title", "旗舰")
	searchResult, err := models.EsClient.Search().
		Index("product").
		Query(query).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	var productList []models.Product
	var product models.Product
	for _, item := range searchResult.Each(reflect.TypeOf(product)) {
		g := item.(models.Product)
		log.Infof("标题： %v\n", g.Title)
		productList = append(productList, g)
	}

	c.Data["json"] = productList
	_ = c.ServeJSON()
}

// FilterQuery 条件筛选查询
func (c *SearchController) FilterQuery() {
	defer func() {
		if r := recover(); r != nil {
			log.Info("Recovered in f", r)
			c.Ctx.WriteString("Query")
		}
	}()

	//筛选
	boolQ := elastic.NewBoolQuery()
	boolQ.Must(elastic.NewMatchQuery("Title", "小米"))
	boolQ.Filter(elastic.NewRangeQuery("Id").Gt(19))
	boolQ.Filter(elastic.NewRangeQuery("Id").Lt(31))
	searchResult, err := models.EsClient.Search().Index("product").Query(boolQ).Do(context.Background())
	if err != nil {
		log.Error(err)
	}
	var product models.Product
	for _, item := range searchResult.Each(reflect.TypeOf(product)) {
		t := item.(models.Product)
		log.Infof("Id:%v 标题：%v\n", t.Id, t.Title)
	}

	c.Ctx.WriteString("filter Query")
}

// ProductList 分页查询
func (c *SearchController) ProductList() {
	c.BaseInit()
	keyword := c.GetString("keyword")
	defer func() {
		if r := recover(); r != nil {
			log.Info("Recovered in f", r)
			c.Ctx.WriteString("ProductList")
		}
	}()

	page, _ := c.GetInt("page")
	if page == 0 {
		page = 1
	}
	pageSize := 5

	query := elastic.NewMatchQuery("Title", keyword)
	searchResult, err := models.EsClient.Search().
		Index("product").
		Query(query).
		Sort("Price", true). //true 升序
		Sort("Id", false). //false 降序
		From((page - 1) * pageSize).Size(pageSize).
		Do(context.Background())
	if err != nil {
		panic(err)
	}

	//查询符合条件的商品的总数
	searchResult2, _ := models.EsClient.Search().
		Index("product").
		Query(query).
		Do(context.Background())

	var productList []models.Product
	var product models.Product
	for _, item := range searchResult.Each(reflect.TypeOf(product)) {
		g := item.(models.Product)
		log.Infof("标题： %v\n", g.Title)
		productList = append(productList, g)
	}
	c.Data["productList"] = productList
	c.Data["totalPages"] = math.Ceil(float64(len(searchResult2.Each(reflect.TypeOf(product)))) / float64(pageSize))
	c.Data["page"] = page
	c.Data["keyword"] = keyword
	c.TplName = "frontend/elasticsearch/list.html"
}
