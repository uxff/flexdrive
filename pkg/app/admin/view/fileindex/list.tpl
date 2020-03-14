{{ define "fileindex/list.tpl" }}

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
            <li><a href="/fileindex/list">文件</a></li>
            <li class="active">文件列表</li>
        </ul>
        <div class="panel panel-default">
            <div class="panel-heading">查询条件</div>
            <div class="panel-body">
                <form id="formSearch" class="form-horizontal" method="GET" action="/fileindex/list">
                    <div class="form-group" >
                        <label class="control-label col-sm-1" for="txt_search_name">文件名称</label>
                        <div class="col-sm-2">
                            <input type="text" class="form-control" id="txt_search_name" name="name" value="{{.reqParam.Name}}">
                        </div>
                        <label class="control-label col-sm-1" for="txt_search_hash">Hash</label>
                        <div class="col-sm-2">
                            <input type="text" class="form-control" id="txt_search_hash" name="fileHash" value="{{.reqParam.FileHash}}">
                        </div>
                        <label class="control-label col-sm-1" for="txt_search_created">时间</label>
                        <div class="col-sm-2">
                            <input type="text" class="form-control" id="txt_search_created_start" name="createStart" value="{{.reqParam.CreateStart}}">
                        </div>
                        <div class="col-sm-2">
                            <input type="text" class="form-control" id="txt_search_created_end" name="createEnd" value="{{.reqParam.CreateEnd}}">
                        </div>
                        <div class="col-sm-1" style="text-align:right;">
                            <button type="submit" class="btn btn-primary">查询</button>
                        </div>
                    </div>
                </form>
            </div>
        </div>       

        <table class="table table-striped table-bordered table-hover">
            <thead>
                <tr class="info">
                    <th>id</th>
                    <th>名称</th>
                    <th>Hash</th>
                    <th>节点</th>
                    <th>创建时间</th>
                    <th>引用量</th>
                    <th>大小(B)</th>
                    <th>占用空间(kB)</th>
                    <th>备注</th>
                    <th>状态</th>
                    <th>操作</th>
                </tr>
            </thead>
            <tbody>
                {{range .list}}
                <tr>
                    <td>{{.Id}}</td>
                    <td>{{.FileName}}</td>
                    <td>{{.FileHash}}</td>
                    <td>{{.NodeId }}, {{.NodeId2}}, {{.NodeId2}}</td>
                    <td>{{.Created }}</td>
                    <td>{{.ReferCount }}</td>
                    <td>{{size4Human .Size }}</td>
                    <td>{{space4Human .Space }}</td>
                    <td>{{.Desc}}</td>
                    <td>{{mgrStatus .Status}}</td>
                    <td>
                        {{if eq .Status 1}}
                        <a href="/fileindex/enable/{{.Id}}/9">禁用</a>
                        {{else}}
                        <a href="/fileindex/enable/{{.Id}}/1">启用</a>
                        {{end}}
                        <a href="/file/{{.FileHash}}/{{.FileName}}">下载</a>
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
