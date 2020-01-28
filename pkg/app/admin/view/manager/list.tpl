{{ define "manager/list.tpl" }}

{{template "common/head.tpl" .}}
{{template "common/partheader.tpl" .}}

<!-- 时间选择器样式表 -->
<link href="https://cdn.bootcss.com/bootstrap-datetimepicker/4.17.47/css/bootstrap-datetimepicker.min.css" rel="stylesheet">
<!-- 时间选择器前置脚本 -->
<script src="https://cdn.bootcss.com/moment.js/2.22.1/moment-with-locales.min.js"></script>
<!-- 时间选择器核心脚本 -->
<script src="https://cdn.bootcss.com/bootstrap-datetimepicker/4.17.47/js/bootstrap-datetimepicker.min.js"></script>

<div class="container">

    <div class="row">
        <ul class="breadcrumb">
            <li><a href="/">首页</a></li>
            <li><a href="/manager/list">管理员</a></li>
            <li class="active">管理员列表</li>
        </ul>
        <div class="panel panel-default">
            <div class="panel-heading">查询条件</div>
            <div class="panel-body">
                <form id="formSearch" class="form-horizontal" method="GET" action="/manager/list">
                    <div class="form-group" >
                        <label class="control-label col-sm-1" for="txt_search_name">Email</label>
                        <div class="col-sm-3">
                            <input type="text" class="form-control" id="txt_search_name" name="email" value="{{.reqParam.Email}}">
                        </div>
                        <label class="control-label col-sm-1" for="txt_search_created">时间</label>
                        <div class="col-sm-2">
                            <input type="text" class="form-control" id="txt_search_created_start" name="createStart" value="{{.reqParam.CreateStart}}">
                        </div>
                        <div class="col-sm-2">
                            <input type="text" class="form-control" id="txt_search_created_end" name="createEnd" value="{{.reqParam.CreateEnd}}">
                        </div>
                        <div class="col-sm-3" style="text-align:left;">
                            <button type="submit" style="margin-left:50px" class="btn btn-primary">查询</button>
                            <a href="/manager/add" type="button" style="margin-left:50px" class="btn btn-success">新增</a>
                        </div>
                    </div>
                </form>
            </div>
        </div>       

        <div id="toolbar" class="btn-group" style="display: none">
            <button id="btn_add" type="button" class="btn btn-default">
                <span class="glyphicon glyphicon-plus" aria-hidden="true"></span>新增
            </button>
            <button id="btn_edit" type="button" class="btn btn-default">
                <span class="glyphicon glyphicon-pencil" aria-hidden="true"></span>修改
            </button>
            <button id="btn_delete" type="button" class="btn btn-default">
                <span class="glyphicon glyphicon-remove" aria-hidden="true"></span>删除
            </button>
        </div>
        <table id="table-manager-list" class="table table-striped table-bordered table-hover">
            <thead>
                <tr class="info">
                    <th>mid</th>
                    <th>Email</th>
                    <th>角色</th>
                    <th>最后登录时间</th>
                    <th>最后登录ip</th>
                    <th>状态</th>
                    <th>操作</th>
                </tr>
            </thead>
            <tbody>
                {{range .list}}
                <tr>
                    <td>{{.Mid}}</td>
                    <td>{{.Email}}</td>
                    <td>{{.RoleName}}</td>
                    <td>{{.LastLoginAt}}</td>
                    <td>{{.LastLoginIp}}</td>
                    <td>{{mgrStatus .Status}}</td>
                    <td>
                        {{if eq .Status 1}}
                        <a href="/manager/enable/{{.Mid}}/9">禁用</a>
                        {{else}}
                        <a href="/manager/enable/{{.Mid}}/1">启用</a>
                        {{end}}
                        <a href="/manager/edit/{{.Mid}}">编辑</a>
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

<script type="text/javascript">
    $("#txt_search_created_start").datetimepicker({
        format: 'YYYY-MM-DD HH:mm'
    });
    $("#txt_search_created_end").datetimepicker({
        format: 'YYYY-MM-DD HH:mm'
    });
</script>
    
{{template "common/partfooter.tpl"}}
{{template "common/foot.tpl"}}

{{ end }}
