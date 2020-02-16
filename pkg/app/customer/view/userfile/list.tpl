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
        <div class="col-md-5" style="padding: 5px;">
            当前等级：黄金会员 &nbsp;&nbsp;
            当前空间：已用 {{space4Human .LoginInfo.UserEnt.UsedSpace}} / 总共 {{space4Human .LoginInfo.UserEnt.QuotaSpace}}
            &nbsp;&nbsp;
            <a href="/" style="text-align: right;" >扩容</a>
            <div class="progress " style="width:100%; float: left; height: 4px; margin-bottom: 10px;">
                <div class="progress-bar progress-bar-success" role="progressbar"
                        aria-valuenow="60" aria-valuemin="0" aria-valuemax="100"
                        style="width: {{spaceRate .LoginInfo.UserEnt.UsedSpace .LoginInfo.UserEnt.QuotaSpace}}%;">
                    <span class="sr-only">{{spaceRate .LoginInfo.UserEnt.UsedSpace .LoginInfo.UserEnt.QuotaSpace}}% 已用</span>
                </div>
            </div>
        </div>
        <div class="col-md-1">
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
            <button class="btn btn-info " type="button" data-toggle="modal" data-target="#newFolderModal">新建文件夹</button>
            <a href="javascript:;" class="btn btn-info " type="button" data-toggle="modal" data-target="#uploadModal"><span class="glyphicon glyphicon-cloud-upload"></span>上传</a>
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
            <a href="/my/file/list?dir={{.parentPath}}"><span class="glyphicon glyphicon-circle-arrow-left"></span>返回上一级</a>&nbsp;
                位置：
            <li><a href="/my/file/list">全部文件</a></li>
            {{range $lk, $lv := .dirLis}}
            {{if $lv}} 
            <li><a href="/my/file/list?dir={{$lv.Parent}}{{$lv.Dir}}">{{$lv.Dir}}</a></li>
            {{end}}
            {{end}}
            <li><input type="hidden" id="dirPath" value="{{.reqParam.Dir}}" readonly></li>
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
                    <td>
                        {{if .IsDir}}<span class="glyphicon glyphicon-folder-close"></span>{{else}}<span class="glyphicon glyphicon-file"></span>{{end}}
                        <a href="/my/file/list?dir={{.FilePath}}{{.FileName}}">{{.FileName}}</a>
                    </td>
                    <td>{{.Created }}</td>
                    <td>{{.Size }}</td>
                    <td>
                        <a href="/">移动到</a>
                        <a href="/">复制到</a>
                        <a href="/">重命名</a>
                        {{if eq .Status 1}}
                        <a href="/my/file/enable/{{.Id}}/9">删除</a>
                        {{end}}
                        <a href="/">分享</a>
                        <a href="/">下载</a>
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

<!-- 模态框（Modal） newFolder -->
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
            <form id="newFolderForm" accept-charset="utf-8" role="form" class="form-horizontal" method="POST" action='/my/file/newfolder'>

            <div class="modal-body">
                <div class="row">
                    <div class="col-md-4 text-right">
                        当前路径：
                    </div>
                    <div class="col-md-6">
                        全部文件<span id="dirPathTextInNewFolderModal"></span>
                        <input type="hidden" name="parentDir" id="dirPathInNewFolderModal" readonly value="{{.reqParam.Dir}}">
                    </div>
                </div>
                <div class="row">
                    <div class="col-md-4 text-right">
                        请输入文件夹名称：
                    </div>
                    <div class="col-md-6">
                        <input type="text" name="dirName" id="nameInNewFolderModal">
                    </div>
                </div>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-default" data-dismiss="modal">关闭</button>
                <button type="submit" class="btn btn-primary" id="newFolderSubmit">提交</button>
            </div>
            </form>
        </div><!-- /.modal-content -->
    </div><!-- /.modal-dialog -->
</div><!-- /.modal -->
    
<!-- 模态框（Modal） upload -->
<div class="modal fade" id="uploadModal" tabindex="-1" role="dialog" aria-labelledby="uploadLabel" aria-hidden="true">
    <div class="modal-dialog">
        <div class="modal-content">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal" 
                        aria-hidden="true">×
                </button>
                <h4 class="modal-title" id="uploadLabel">
                    上传文件
                </h4>
            </div>
            <form id="uploadForm" accept-charset="utf-8" role="form" class="form-horizontal" method="POST" action='/my/file/upload' enctype="multipart/form-data">
            <div class="modal-body">
                <div class="row">
                    <div class="col-md-4 text-right">
                        当前路径：
                    </div>
                    <div class="col-md-6">
                        全部文件<span id="dirPathTextInUploadModal"></span>
                        <input type="hidden" name="parentDir" id="dirPathInUploadModal" readonly value="{{.reqParam.Dir}}">
                    </div>
                </div>
                <div class="row">
                    <div class="col-md-4 text-right">
                        请选择文件：
                    </div>
                    <div class="col-md-6">
                        <input type="file" name="file" id="fileInUploadModal">
                    </div>
                </div>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-default" data-dismiss="modal">关闭</button>
                <button type="button" class="btn btn-primary" id="uploadSubmit">提交</button>
            </div>
            </form>
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

$(function () {
    $('#newFolderModal').on('show.bs.modal', function () {
        //alert('嘿，我听说您喜欢模态框xxxxxxxxx...');})
        $('#dirPathTextInNewFolderModal').html($('#dirPath').val());
        $('#nameInNewFolderModal').focus();// 未生效
    });
    $('#newFolderSubmit').on('click', function(){
        $('#newFolderForm').submit();
        $('#newFolderModal').modal('hide');
    });
    $('#uploadModal').on('show.bs.modal', function () {
        //alert('嘿，我听说您喜欢模态框xxxxxxxxx...');})
        $('#dirPathTextInUploadModal').html($('#dirPath').val());
    });
    $('#uploadSubmit').on('click', function(){
        $('#uploadForm').submit();
        $('#uploadModal').modal('hide');
    });
});

</script>
    
{{template "common/partfooter.tpl"}}
{{template "common/foot.tpl"}}

{{ end }}
