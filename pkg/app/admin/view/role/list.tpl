{{ define "role/list.tpl" }}

{{template "common/head.tpl" .}}
{{template "common/partheader.tpl" .}}


<div class="container">

    <div class="row">
        <ul class="breadcrumb">
            <li><a href="/">首页</a></li>
            <li><a href="/role/list">角色</a></li>
            <li class="active">角色列表</li>
        </ul>
        <div class="panel panel-default">
            <div class="panel-heading">查询条件</div>
            <div class="panel-body">
                <form id="formSearch" class="form-horizontal" method="GET" action="/role/list">
                    <div class="form-group" >
                        <label class="control-label col-sm-1" for="txt_search_departmentname">名称</label>
                        <div class="col-sm-3">
                            <input type="text" class="form-control" id="txt_search_departmentname">
                        </div>
                        <label class="control-label col-sm-1" for="txt_search_statu">时间</label>
                        <div class="col-sm-3">
                            <input type="text" class="form-control" id="txt_search_statu">
                        </div>
                        <div class="col-sm-4" style="text-align:left;">
                            <button type="button" style="margin-left:50px" class="btn btn-primary">查询</button>
                            <a href="/role/add" type="button" style="margin-left:50px" class="btn btn-success">新增</a>
                        </div>
                    </div>
                </form>
            </div>
        </div>       

        <table id="table-manager-list" class="table table-striped table-bordered table-hover">
            <thead>
                <tr class="info">
                    <th>id</th>
                    <th>名称</th>
                    <th>状态</th>
                    <th>操作</th>
                </tr>
            </thead>
            <tbody>
                {{range .list}}
                <tr>
                    <td>{{.Id}}</td>
                    <td>{{.Name}}</td>
                    <td>{{mgrStatus .Status}}</td>
                    <td>
                        {{if eq .Status 1}}
                        <a href="/role/enable/{{.Mid}}/9">禁用</a>
                        {{else}}
                        <a href="/role/enable/{{.Mid}}/1">启用</a>
                        {{end}}
                        <a href="/role/edit/{{.Mid}}">编辑</a>
                    </td>
                </tr>
                {{else}}
                <tr>
                    <td colspan="12">暂无数据</td>
                </tr>
                {{end}}
            </tbody>
        </table>
        <!--分页-->
        {{template "paginator2.tpl" .}}


    </div>

</div>

{{template "common/partfooter.tpl"}}
{{template "common/foot.tpl"}}

{{ end }}
