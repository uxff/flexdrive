{{ define "userlevel/list.tpl" }}

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
            <li><a href="/userlevel/list">会员升级包</a></li>
            <li class="active">会员升级包列表</li>
        </ul>
        <div class="panel panel-default">
            <div class="panel-heading">查询条件</div>
            <div class="panel-body">
                <form id="formSearch" class="form-horizontal" method="GET" action="/userlevel/list">
                    <div class="form-group" >
                        <label class="control-label col-sm-1" for="txt_search_name">名称</label>
                        <div class="col-sm-2">
                            <input type="text" class="form-control" id="txt_search_name" name="name" value="{{.reqParam.Name}}">
                        </div>
                        <label class="control-label col-sm-1" for="txt_search_created">创建时间</label>
                        <div class="col-sm-2">
                            <input type="text" class="form-control" id="txt_search_created_start" name="createStart" value="{{.reqParam.CreateStart}}">
                        </div>
                        <div class="col-sm-2">
                            <input type="text" class="form-control" id="txt_search_created_end" name="createEnd" value="{{.reqParam.CreateEnd}}">
                        </div>
                        <div class="col-sm-3" style="text-align:left;">
                            <button type="submit" style="margin-left:50px" class="btn btn-primary">查询</button>
                            <a href="/userlevel/add" type="button" style="margin-left:50px" class="btn btn-success">新增</a>
                        </div>
                    </div>
                </form>
            </div>
        </div>       

        <table id="table-list" class="table table-striped table-bordered table-hover">
            <thead>
                <tr class="info">
                    <th>id</th>
                    <th>升级包名称</th>
                    <th>介绍</th>
                    <th>是否默认</th>
                    <th>配额空间</th>
                    <th>价格(元)</th>
                    <th>创建时间</th>
                    <th>状态</th>
                    <th>操作</th>
                </tr>
            </thead>
            <tbody>
                {{range .list}}
                <tr>
                    <td>{{.Id}}</td>
                    <td>{{.Name}}</td>
                    <td>{{.Desc}}</td>
                    <td>{{if eq .IsDefault 1}}是{{else}}否{{end}}</td>
                    <td>{{space4Human .QuotaSpace}}</td>
                    <td>{{amount4Human .Price}}</td>
                    <td>{{.Created}}</td>
                    <td>{{mgrStatus .Status}}</td>
                    <td>
                        {{if eq .Status 1}}
                        <a href="/userlevel/enable/{{.Id}}/9">禁用</a>
                        {{else}}
                        <a href="/userlevel/enable/{{.Id}}/1">启用</a>
                        {{end}}
                        <a href="/userlevel/edit/{{.Id}}">编辑</a>
                    </td>
                </tr>
                {{else}}
                <tr>
                    <td colspan="12" class="text-center">暂无数据</td>
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
