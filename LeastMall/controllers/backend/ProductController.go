package backend

import (
	"LeastMall/common"
	"LeastMall/models"
	"fmt"
	"github.com/prometheus/common/log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
)

var wg sync.WaitGroup

type ProductController struct {
	BaseController
}

func (c *ProductController) Get() {
	page, err := c.GetInt("page")
	if err != nil || page == 0 {
		page = 1
	}
	keyword := c.GetString("keyword")
	pageSize := 5

	where := "1=1" // 1=1 作用相当于select * from table；这里为拼接字符串方便
	if len(keyword) > 0 {
		where += " AND title like \"%" + keyword + "%\""
	}
	var productList []models.Product
	models.DB.Where(where).Offset((page - 1) * pageSize).Limit(pageSize).Find(&productList)
	var count int
	models.DB.Where(where).Table("product").Count(&count)
	c.Data["productList"] = productList
	c.Data["totalPages"] = math.Ceil(float64(count) / float64(pageSize))
	c.Data["page"] = page
	c.TplName = "backend/product/index.html"
}

func (c *ProductController) Add() {
	// 获取商品分类
	var productCateList []models.ProductCate
	models.DB.Where("pid=?", 0).Preload("ProductCateItem").Find(&productCateList)
	c.Data["productCateList"] = productCateList

	// 获取颜色信息
	var productColor []models.ProductColor
	models.DB.Find(&productColor)
	c.Data["productColor"] = productColor

	// 获取商品类型信息
	var productType []models.ProductType
	models.DB.Find(&productType)
	c.Data["productType"] = productType

	c.TplName = "backend/product/add.html"
}

func (c *ProductController) GoAdd() {
	// 1、获取表单提交过来的数据
	title := c.GetString("title")
	subTitle := c.GetString("sub_title")
	productSn := c.GetString("product_sn")
	cateId, _ := c.GetInt("cate_id")
	productNumber, _ := c.GetInt("product_number")
	marketPrice, _ := c.GetFloat("market_price")
	price, _ := c.GetFloat("price")
	relationProduct := c.GetString("relation_product")
	productAttr := c.GetString("product_attr")
	productVersion := c.GetString("product_version")
	productGift := c.GetString("product_gift")
	productFitting := c.GetString("product_fitting")
	productColor := c.GetStrings("product_color")
	productKeywords := c.GetString("product_keywords")
	productDesc := c.GetString("product_desc")
	productContent := c.GetString("product_content")
	isDelete, _ := c.GetInt("is_delete")
	isHot, _ := c.GetInt("is_hot")
	isBest, _ := c.GetInt("is_best")
	isNew, _ := c.GetInt("is_new")
	productTypeId, _ := c.GetInt("product_type_id")
	sort, _ := c.GetInt("sort")
	status, _ := c.GetInt("status")
	addTime := int(common.GetUnix())

	// 2、获取颜色信息 把颜色转化成字符串
	productColorStr := strings.Join(productColor, ",")

	// 3、上传图片   生成缩略图
	productImg, _ := c.UploadImg("product_img")

	// 4、增加商品数据
	product := models.Product{
		Title:           title,
		SubTitle:        subTitle,
		ProductSn:       productSn,
		CateId:          cateId,
		ClickCount:      100,
		ProductNumber:   productNumber,
		MarketPrice:     marketPrice,
		Price:           price,
		RelationProduct: relationProduct,
		ProductAttr:     productAttr,
		ProductVersion:  productVersion,
		ProductGift:     productGift,
		ProductFitting:  productFitting,
		ProductKeywords: productKeywords,
		ProductDesc:     productDesc,
		ProductContent:  productContent,
		IsDelete:        isDelete,
		IsHot:           isHot,
		IsBest:          isBest,
		IsNew:           isNew,
		ProductTypeId:   productTypeId,
		Sort:            sort,
		Status:          status,
		AddTime:         addTime,
		ProductColor:    productColorStr,
		ProductImg:      productImg,
	}
	err1 := models.DB.Create(&product).Error
	if err1 != nil {
		c.Error("增加失败", "/product/add")
	}

	//5、增加图库 信息
	wg.Add(1)
	go func() {
		productImageList := c.GetStrings("product_image_list")
		for _, s := range productImageList {
			productImgObj := models.ProductImage{}
			productImgObj.ProductId = product.Id
			productImgObj.ImgUrl = s
			productImgObj.Sort = 10
			productImgObj.Status = 1
			productImgObj.AddTime = int(common.GetUnix())
			models.DB.Create(&productImgObj)
		}
		wg.Done()
	}()

	// 6、增加规格包装
	wg.Add(1)
	go func() {
		attrIdList := c.GetStrings("attr_id_list")
		attrValueList := c.GetStrings("attr_value_list")
		for i := 0; i < len(attrIdList); i++ {
			productTypeAttributeId, _ := strconv.Atoi(attrIdList[i])
			productTypeAttributeObj := models.ProductTypeAttribute{Id: productTypeAttributeId}
			models.DB.Find(&productTypeAttributeObj)

			productAttrObj := models.ProductAttr{}
			productAttrObj.ProductId = product.Id
			productAttrObj.AttributeTitle = productTypeAttributeObj.Title
			productAttrObj.AttributeType = productTypeAttributeObj.AttrType
			productAttrObj.AttributeId = productTypeAttributeObj.Id
			productAttrObj.AttributeCateId = productTypeAttributeObj.CateId
			productAttrObj.AttributeValue = attrValueList[i]
			productAttrObj.Status = 1
			productAttrObj.Sort = 10
			productAttrObj.AddTime = int(common.GetUnix())
			models.DB.Create(&productAttrObj)
		}
		wg.Done()
	}()

	wg.Wait()
	c.Success("增加数据成功", "/product")
}

func (c *ProductController) Edit() {

	// 1、获取商品数据
	id, err1 := c.GetInt("id")
	if err1 != nil {
		c.Error("非法请求", "/product")
	}
	product := models.Product{Id: id}
	models.DB.Find(&product)
	c.Data["product"] = product

	//2、获取商品分类
	var productCate []models.ProductCate
	models.DB.Where("pid=?", 0).Preload("ProductCateItem").Find(&productCate)
	c.Data["productCateList"] = productCate

	// 3、获取所有颜色 以及选中的颜色
	productColorSlice := strings.Split(product.ProductColor, ",")
	productColorMap := make(map[string]string)
	for _, v := range productColorSlice {
		productColorMap[v] = v
	}
	//获取颜色信息
	var productColor []models.ProductColor
	models.DB.Find(&productColor)
	for i := 0; i < len(productColor); i++ {
		_, ok := productColorMap[strconv.Itoa(productColor[i].Id)]
		if ok {
			productColor[i].Checked = true
		}
	}
	c.Data["productColor"] = productColor

	//4、商品的图库信息
	var productImage []models.ProductImage
	models.DB.Where("product_id=?", product.Id).Find(&productImage)
	c.Data["productImage"] = productImage

	// 5、获取商品类型
	var productType []models.ProductType
	models.DB.Find(&productType)
	c.Data["productType"] = productType
	//6、获取规格信息
	var productAttr []models.ProductAttr
	models.DB.Where("product_id=?", product.Id).Find(&productAttr)

	var productAttrStr string
	for _, v := range productAttr {
		if v.AttributeType == 1 {
			productAttrStr += fmt.Sprintf(`<li><span>%v: 　</span>  <input type="hidden" name="attr_id_list" value="%v" />   <input type="text" name="attr_value_list" value="%v" /></li>`, v.AttributeTitle, v.AttributeId, v.AttributeValue)
		} else if v.AttributeType == 2 {
			productAttrStr += fmt.Sprintf(`<li><span>%v: 　</span><input type="hidden" name="attr_id_list" value="%v" />  <textarea cols="50" rows="3" name="attr_value_list">%v</textarea></li>`, v.AttributeTitle, v.AttributeId, v.AttributeValue)
		} else {

			// 获取 attr_value  获取可选值列表
			oneProductTypeAttribute := models.ProductTypeAttribute{Id: v.AttributeId}
			models.DB.Find(&oneProductTypeAttribute)
			attrValueSlice := strings.Split(oneProductTypeAttribute.AttrValue, "\n")
			productAttrStr += fmt.Sprintf(`<li><span>%v: 　</span>  <input type="hidden" name="attr_id_list" value="%v" /> `, v.AttributeTitle, v.AttributeId)
			productAttrStr += fmt.Sprintf(`<select name="attr_value_list">`)
			for j := 0; j < len(attrValueSlice); j++ {
				if attrValueSlice[j] == v.AttributeValue {
					productAttrStr += fmt.Sprintf(`<option value="%v" selected >%v</option>`, attrValueSlice[j], attrValueSlice[j])
				} else {
					productAttrStr += fmt.Sprintf(`<option value="%v">%v</option>`, attrValueSlice[j], attrValueSlice[j])
				}
			}
			productAttrStr += fmt.Sprintf(`</select>`)
			productAttrStr += fmt.Sprintf(`</li>`)
		}
	}

	c.Data["productAttrStr"] = productAttrStr
	//上一页地址
	c.Data["prevPage"] = c.Ctx.Request.Referer()
	c.TplName = "backend/product/edit.html"
}

func (c *ProductController) GoEdit() {

	//1.获取要修改的商品数据
	id, err1 := c.GetInt("id")
	if err1 != nil {
		c.Error("非法请求", "/product")
	}
	title := c.GetString("title")
	subTitle := c.GetString("sub_title")
	productSn := c.GetString("product_sn")
	cateId, _ := c.GetInt("cate_id")
	productNumber, _ := c.GetInt("product_number")
	marketPrice, _ := c.GetFloat("market_price")
	price, _ := c.GetFloat("price")
	relationProduct := c.GetString("relation_product")
	productAttr := c.GetString("product_attr")
	productVersion := c.GetString("product_version")
	productGift := c.GetString("product_gift")
	productFitting := c.GetString("product_fitting")
	productColor := c.GetStrings("product_color")
	productKeywords := c.GetString("product_keywords")
	productDesc := c.GetString("product_desc")
	productContent := c.GetString("product_content")
	isDelete, _ := c.GetInt("is_delete")
	isHot, _ := c.GetInt("is_hot")
	isBest, _ := c.GetInt("is_best")
	isNew, _ := c.GetInt("is_new")
	productTypeId, _ := c.GetInt("product_type_id")
	sort, _ := c.GetInt("sort")
	status, _ := c.GetInt("status")

	prevPage := c.GetString("prevPage")
	//2.获取颜色信息 把颜色转化成字符串
	productColorStr := strings.Join(productColor, ",")

	product := models.Product{Id: id}
	models.DB.Find(&product)
	product.Title = title
	product.SubTitle = subTitle
	product.ProductSn = productSn
	product.CateId = cateId
	product.ProductNumber = productNumber
	product.MarketPrice = marketPrice
	product.Price = price
	product.RelationProduct = relationProduct
	product.ProductAttr = productAttr
	product.ProductVersion = productVersion
	product.ProductGift = productGift
	product.ProductFitting = productFitting
	product.ProductKeywords = productKeywords
	product.ProductDesc = productDesc
	product.ProductContent = productContent
	product.IsDelete = isDelete
	product.IsHot = isHot
	product.IsBest = isBest
	product.IsNew = isNew
	product.ProductTypeId = productTypeId
	product.Sort = sort
	product.Status = status
	product.ProductColor = productColorStr

	//3.上传图片，生成缩略图
	productImg, err2 := c.UploadImg("product_img")
	log.Info("img-->", productImg)
	if err2 == nil && len(productImg) > 0 {
		product.ProductImg = productImg
	}
	//4.执行修改商品
	err3 := models.DB.Save(&product).Error
	if err3 != nil {
		c.Error("修改数据失败", "/product/edit?id="+strconv.Itoa(id))
		return
	}
	//5.修改图库数据 （增加）
	wg.Add(1)
	go func() {
		productImageList := c.GetStrings("product_image_list")
		for _, v := range productImageList {
			productImgObj := models.ProductImage{}
			productImgObj.ProductId = product.Id
			productImgObj.ImgUrl = v
			productImgObj.Sort = 10
			productImgObj.Status = 1
			productImgObj.AddTime = int(common.GetUnix())
			models.DB.Create(&productImgObj)
		}
		wg.Done()
	}()

	//6.修改商品类型属性数据         1、删除当前商品id对应的类型属性  2、执行增加

	//删除当前商品id对应的类型属性
	productAttrObj := models.ProductAttr{}
	models.DB.Where("product_id=?", product.Id).Delete(&productAttrObj)
	//执行增加
	wg.Add(1)
	go func() {
		attrIdList := c.GetStrings("attr_id_list")
		attrValueList := c.GetStrings("attr_value_list")
		for i := 0; i < len(attrIdList); i++ {
			productTypeAttributeId, _ := strconv.Atoi(attrIdList[i])
			productTypeAttributeObj := models.ProductTypeAttribute{Id: productTypeAttributeId}
			models.DB.Find(&productTypeAttributeObj)

			productAttrObj := models.ProductAttr{}
			productAttrObj.ProductId = product.Id
			productAttrObj.AttributeTitle = productTypeAttributeObj.Title
			productAttrObj.AttributeType = productTypeAttributeObj.AttrType
			productAttrObj.AttributeId = productTypeAttributeObj.Id
			productAttrObj.AttributeCateId = productTypeAttributeObj.CateId
			productAttrObj.AttributeValue = attrValueList[i]
			productAttrObj.Status = 1
			productAttrObj.Sort = 10
			productAttrObj.AddTime = int(common.GetUnix())
			models.DB.Create(&productAttrObj)
		}
		wg.Done()
	}()

	wg.Wait()
	c.Success("修改数据成功", prevPage)
}
func (c *ProductController) Delete() {
	id, err1 := c.GetInt("id")
	if err1 != nil {
		c.Error("传入参数错误", "/product")
		return
	}
	product := models.Product{Id: id}
	models.DB.Find(&product)
	path, _ := os.Getwd()
	path = strings.ReplaceAll(path, "\\", "/")
	address := path + "/" + product.ProductImg
	_ = os.Remove(address)
	err2 := models.DB.Delete(&product).Error
	if err2 != nil {
		productAttr := models.ProductAttr{ProductId: id}
		models.DB.Delete(&productAttr)
		productImage := models.ProductImage{ProductId: id}
		models.DB.Delete(&productImage)
	}
	c.Success("删除商品成功", c.Ctx.Request.Referer())
}

func (c *ProductController) GoUpload() {

	savePath, err := c.UploadImg("file")
	if err != nil {
		log.Error("失败")
		c.Data["json"] = map[string]interface{}{
			"link": "",
		}
		_ = c.ServeJSON()
	} else {
		//返回json数据
		c.Data["json"] = map[string]interface{}{
			"link": "/" + savePath,
		}
		_ = c.ServeJSON()
	}

}

// GetProductTypeAttribute 获取商品类型属性
func (c *ProductController) GetProductTypeAttribute() {
	cateId, err1 := c.GetInt("cate_id")
	var ProductTypeAttribute []models.ProductTypeAttribute
	err2 := models.DB.Where("cate_id=?", cateId).Find(&ProductTypeAttribute).Error
	if err1 != nil || err2 != nil {
		c.Data["json"] = map[string]interface{}{
			"result":  "",
			"success": false,
		}
		_ = c.ServeJSON()

	} else {
		c.Data["json"] = map[string]interface{}{
			"result":  ProductTypeAttribute,
			"success": true,
		}
		_ = c.ServeJSON()
	}
}

// ChangeProductImageColor 修改图片对应颜色信息
func (c *ProductController) ChangeProductImageColor() {
	colorId, err1 := c.GetInt("color_id")
	productImageId, err2 := c.GetInt("product_image_id")
	productImage := models.ProductImage{Id: productImageId}
	models.DB.Find(&productImage)
	productImage.ColorId = colorId
	err3 := models.DB.Save(&productImage).Error

	if err1 != nil || err2 != nil || err3 != nil {
		c.Data["json"] = map[string]interface{}{
			"result":  "更新失败",
			"success": false,
		}
		_ = c.ServeJSON()
	} else {
		c.Data["json"] = map[string]interface{}{
			"result":  "更新成功",
			"success": true,
		}
		_ = c.ServeJSON()
	}
}

// RemoveProductImage 删除图库
func (c *ProductController) RemoveProductImage() {
	productImageId, err1 := c.GetInt("product_image_id")
	productImage := models.ProductImage{Id: productImageId}
	err2 := models.DB.Delete(&productImage).Error
	_ = os.Remove(productImage.ImgUrl)

	if err1 != nil || err2 != nil {
		c.Data["json"] = map[string]interface{}{
			"result":  "删除失败",
			"success": false,
		}
		_ = c.ServeJSON()
	} else {
		//删除图片
		c.Data["json"] = map[string]interface{}{
			"result":  "删除",
			"success": true,
		}
		_ = c.ServeJSON()
	}

}
