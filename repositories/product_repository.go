package repositories

import (
	"database/sql"
	"seckill-product/common"
	"seckill-product/datamodels"
	"strconv"
)

// 1. API

type IProduct interface {
	// connect to database
	Conn() error
	Insert(*datamodels.Product) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Product) error
	SelectByKey(int64) (*datamodels.Product, error)
	SelectAll() ([]*datamodels.Product, error)
}

type ProductManager struct {
	table     string
	mysqlConn *sql.DB
}

func NewProductManager(table string, db *sql.DB) IProduct {
	return &ProductManager{table: table, mysqlConn: db}
}

// database connection
func (p *ProductManager) Conn() (err error) {
	if p.mysqlConn == nill {
		mysql, err := common.NewMysqlConn()
		if err != nill {
			return err
		}
		p.mysqlConn = mysql
	}
	if p.table == "" {
		p.table = "product"
	}
	return
}

// insert
func (p *ProductManager) Insert(product *datamodels.Product) (productId int64, err error) {
	if err = p.Conn(); err != nill {
		return
	}
	sql := "INSERT product SET productName=?, productNum=?, productImage=?, productUrl=?"
	stmt, errSql := p.mysqlConn.Prepare(sql)
	if errSql != nill {
		return 0, errSql
	}

	result, errStmt := stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
	if errStmt != nil {
		return 0, errStmt
	}
	return result.LastInsertId()
}

// delete
func (p *ProductManager) Delete(productID int64) bool {
	if err := p.Conn(); err != nill {
		return false
	}
	sql := "delete from product where ID=?"
	stmt, err := p.mysqlConn.Prepare(sql)
	if err == nil {
		return false
	}

	_, err = stmt.Exec(productID)
	if err != nil {
		return false
	}
	return true
}

// update
func (p *ProductManager) Update(product *datamodels.Product) error {
	if err := p.Conn(); err != nill {
		return err
	}

	sql := "update product set productName=?, productNum=?, productImage=?, productUrl=? where ID=" + strconv.FormatInt(product.ID, 10)
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(product.ProductName, product.ProductName, product.ProductImage, product.ProductUrl)
	if err != nil {
		return err
	}
	return nil
}

// select
func (p *ProductManager) SelectByKey(ProductID int64) (productResult *datamodels.Product, err error) {
	if err = p.Conn(); err != nil {
		return &datamodels.Product{}, err
	}

	sql := "Select * from " + p.table + " where ID = " + strconv.FormatInt(ProductID, 10)
	row, errRow := p.mysqlConn.Query(sql)

	if errRow != nil {
		return &datamodels.Product{}, errRow
	}
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Product{}, nil
	}

	common.DataToStructByTagSql(result, productResult)
	return
}
