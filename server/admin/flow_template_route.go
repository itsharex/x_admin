package admin

import (
	"x_admin/admin/flow_template"
	"x_admin/middleware"

	"github.com/gin-gonic/gin"
)

/**
集成
1. 导入
- 请先提交git避免文件覆盖!!!
- 下载并解压压缩包后，直接复制server、admin文件夹到项目根目录即可

2. 注册路由
请在 admin/entry.go 文件引入FlowTemplateRoute注册路由

3. 后台手动添加菜单和按钮
flow_template:add
flow_template:edit
flow_template:del
flow_template:list
flow_template:listAll
flow_template:detail
*/

// FlowTemplateRoute(rg)
func FlowTemplateRoute(rg *gin.RouterGroup) {

	handle := flow_template.FlowTemplateHandler{}

	rg = rg.Group("/", middleware.TokenAuth())
	rg.GET("/flow_template/list", handle.List)
	rg.GET("/flow_template/listAll", handle.ListAll)
	rg.GET("/flow_template/detail", handle.Detail)
	rg.POST("/flow_template/add", handle.Add)
	rg.POST("/flow_template/edit", handle.Edit)
	rg.POST("/flow_template/del", handle.Del)
}
