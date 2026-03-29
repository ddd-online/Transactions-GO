package workspace

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/billadm/models"
)

var defaultData = map[string]map[string][]string{
	"expense": {
		"餐饮美食": {"三餐", "商场", "外卖", "奶茶", "零食", "水果", "咖啡", "饮料", "茶叶", "买菜"},
		"交通出行": {"打车", "地铁", "公交", "高铁", "油费", "停车", "ETC", "车险"},
		"购物消费": {"衣物", "数码", "家居", "书籍", "礼物", "玩具", "宠物", "游戏", "快递", "彩票", "电影", "运动", "酒店", "烟酒", "充值", "汽车", "还款"},
		"娱乐休闲": {},
		"生活缴费": {"房租", "物业", "水电", "燃气", "通讯", "人险", "还款", "网费", "理发"},
		"医疗健康": {"医药", "医险"},
		"人情往来": {"红包", "请客", "礼金"},
		"教育学习": {},
	},
	"income": {
		"工资奖金": {"工资", "奖金"},
		"补贴补助": {},
		"退税退款": {},
		"二手转卖": {},
		"彩票收入": {},
		"投资理财": {},
		"借贷借款": {},
		"红包转账": {},
	},
	"transfer": {
		"五险一金": {"养老", "医疗", "失业", "住房"},
		"税费党费": {"团费", "交税"},
	},
}

func seedData(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.Category{}).Count(&count).Error; err != nil {
		logrus.Errorf("检查分类数据失败: %v", err)
		return err
	}
	if count > 0 {
		logrus.Info("数据库已存在数据，跳过预置")
		return nil
	}

	for transactionType, categories := range defaultData {
		for categoryName, tags := range categories {
			category := models.Category{
				Name:            categoryName,
				TransactionType: transactionType,
			}
			if err := db.FirstOrCreate(&category, category).Error; err != nil {
				logrus.Errorf("创建分类失败: %v", err)
				return err
			}

			categoryTransactionType := categoryName + ":" + transactionType
			for _, tagName := range tags {
				tag := models.Tag{
					Name:                     tagName,
					CategoryTransactionType: categoryTransactionType,
				}
				if err := db.FirstOrCreate(&tag, tag).Error; err != nil {
					logrus.Errorf("创建标签失败: %v", err)
					return err
				}
			}
		}
	}
	return nil
}
