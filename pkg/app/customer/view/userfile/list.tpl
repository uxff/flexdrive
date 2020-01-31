{{ define "userfile/list.tpl" }}

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
            <li><a href="/my/file/list">我的文件</a></li>
            <li class="active">文件列表</li>
        </ul>
    </div>

    <!--当前排版方式1-->
    <div class="row" style="margin-bottom: 4px;">
        <div class="col-md-3" style="padding: 5px;">
            我的空间：已用{{kBforHuman .LoginInfo.UserEnt.UsedSpace}} / 总共{{kBforHuman .LoginInfo.UserEnt.QuotaSpace}}
        </div>
        <div class="col-md-3">
            <div class="progress " style="width:100%; margin: 5px; float: left;">
                <div class="progress-bar progress-bar-success" role="progressbar"
                        aria-valuenow="60" aria-valuemin="0" aria-valuemax="100"
                        style="width: 0%;">
                    <span class="sr-only">0% 已用</span>
                </div>
            </div>
        </div>
        <div class="col-md-3">
            <form class="form-horizontal" method="GET" action="/my/file/list">
                <div class="input-group " style="margin-top:0px; position:relative">
                    <input type="text" class="form-control" name="fileName" placeholder="搜索我的文件" value="{{.reqParam.FileName}}" />
                    <span class="input-group-btn">
                        <button class="btn btn-info btn-search" type="submit">🔍搜索</button>
                    </span>
                </div>
            </form>
        </div>
        <div class="col-md-3">
            <button class="btn btn-info " type="button">新建文件夹</button>
            <a href="javascript:;" class="btn btn-info " type="button"><span class="glyphicon glyphicon-cloud-upload"></span>上传</a>
            <button class="btn btn-info " type="button">离线下载</button>
        </div>
    </div>

    <!--保留排版方式2-->
    <div class="row hidden">
        <nav class="navbar navbar-default" role="navigation" style="margin: 0px;"> 
            <div class="container-fluid"> 
                <div class="navbar-header">
                    <button class="btn btn-info navbar-btn" type="button">新建文件夹</button>
                    <a href="javascript:;" class="btn btn-info navbar-btn" type="button"><span class="glyphicon glyphicon-cloud-upload"></span>上传</a>
                    <button class="btn btn-info navbar-btn" type="button">离线下载</button>
                </div> 

                <div>
                    <form class="navbar-form navbar-right" role="search" method="GET" action="/my/file/list">
                        <div class="form-group">
                            <input type="text" class="form-control" name="fileName" value="{{.reqParam.FileName}}" placeholder="搜索我的文件">
                        </div>
                        <button type="submit" class="btn btn-default">🔍搜索</button>
                    </form>
                </div>
            </div> 
        </nav>
    </div>

    <div class="row">
        <ul class="breadcrumb" style="margin: 0px;">
            位置：
            <li><a href="/">全部文件</a></li>
            <li><a href="/my/file/list">文件</a></li>
            <li class="active">文件列表</li>
        </ul>
    </div>
    <div class="row">

        <table class="table table-striped table-bordered table-hover">
            <thead>
                <tr class="info">
                    <th>名称</th>
                    <th>创建时间</th>
                    <th>大小(B)</th>
                    <th>操作</th>
                </tr>
            </thead>
            <tbody>
                {{range .list}}
                <tr data-id="{{.Id}}" data-hash="{{.FileHash}}" data-pathhash={{.PathHash}}>
                    <td>{{.FileName}}</td>
                    <td>{{.Created }}</td>
                    <td>{{.Size }}</td>
                    <td>
                        {{if eq .Status 1}}
                        <a href="/fileindex/enable/{{.Id}}/9">删除</a>
                        {{end}}
                        <a href="/">分享</a>
                        <a href="/">下载</a>
                        <a href="/">移动到</a>
                        <a href="/">复制到</a>
                        <a href="/">重命名</a>
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

<!-- 模态框（Modal） -->
<div class="modal fade" id="newFolderModal" tabindex="-1" role="dialog" aria-labelledby="newFolderLabel" aria-hidden="true">
    <div class="modal-dialog">
        <div class="modal-content">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal" 
                        aria-hidden="true">×
                </button>
                <h4 class="modal-title" id="newFolderLabel">
                    新建文件夹
                </h4>
            </div>
            <div class="modal-body">
                当前路径：
                <input type="">
                请输入文件夹名称：
                <input type="">
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-default" data-dismiss="modal">关闭</button>
                <button type="button" class="btn btn-primary" data-dismiss="modal">提交</button>
            </div>
        </div><!-- /.modal-content -->
    </div><!-- /.modal-dialog -->
</div><!-- /.modal -->
    

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
