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
        <div class="col-md-4" style="padding: 5px;">
            <!--当前等级：{{.userLevel.Name}} &nbsp;&nbsp;-->
            当前空间：已用 {{space4Human .LoginInfo.UserEnt.UsedSpace}} / 总共 {{space4Human .LoginInfo.UserEnt.QuotaSpace}}
            [<a href="/my/order/create" style="text-align: right;" >升级扩容</a>]
            <div class="progress " style="width:100%; float: left; height: 6px; margin-bottom: 10px; background-color: #dff0d8;">
                <div class="progress-bar progress-bar-success" role="progressbar"
                        aria-valuenow="60" aria-valuemin="0" aria-valuemax="100"
                        style="width: {{spaceRate .LoginInfo.UserEnt.UsedSpace .LoginInfo.UserEnt.QuotaSpace}}%;">
                    <span class="sr-only">{{spaceRate .LoginInfo.UserEnt.UsedSpace .LoginInfo.UserEnt.QuotaSpace}}% 已用</span>
                </div>
            </div>
        </div>
        <div class="col-md-4">
            <form class="form-horizontal" method="GET" action="/my/file/list">
                <div class="input-group " style="margin-top:0px; position:relative">
                    <input type="text" class="form-control" name="fileName" placeholder="搜索我的文件" value="{{.reqParam.FileName}}" />
                    <span class="input-group-btn">
                        <button class="btn btn-info btn-search" type="submit">🔍搜索</button>
                    </span>
                </div>
            </form>
        </div>
        <div class="col-md-4">
            <button class="btn btn-info " type="button" data-toggle="modal" data-target="#newFolderModal">新建文件夹</button>
            <a href="javascript:;" class="btn btn-info " type="button" data-toggle="modal" data-target="#uploadModal"><span class="glyphicon glyphicon-cloud-upload"></span>上传</a>
            <button class="btn btn-info " type="button" data-toggle="modal" data-target="#offlineTaskModal">离线下载</button>
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
            <a href="/my/file/list?dir={{.parentPath}}"><span class="glyphicon glyphicon-circle-arrow-up"></span>返回上一级</a>&nbsp;
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
                    <th>大小</th>
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
                    <td>{{if .IsDir}}-{{else}}{{size4Human .Size }}{{end}}</td>
                    <td>
                        <!--
                        <a href="/">移动到</a>
                        <a href="/">复制到</a>
                        -->
                        <a href="#" onclick="checkFile({{.Id}},'{{.FileName}}');"  data-id="{{.Id}}" data-toggle="modal" data-target="#renameModal">重命名</a>
                        {{if eq .Status 1}}
                        <a href="/my/file/enable/{{.Id}}/9">删除</a>
                        {{end}}
                        {{if eq .IsDir 0}}
                        <a href="#" onclick="checkFile({{.Id}},'{{.FileName}}');" data-id="{{.Id}}" data-toggle="modal" data-target="#shareModal">分享</a>
                        <a href="/file/{{.FileHash}}/{{.FileName}}">下载</a>
                        {{end}}
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
                <div class="row" style="margin: 10px;">
                    <div class="col-md-4 text-right">
                        当前路径：
                    </div>
                    <div class="col-md-6">
                        全部文件<span id="dirPathTextInNewFolderModal"></span>
                        <input type="hidden" name="parentDir" id="dirPathInNewFolderModal" readonly value="{{.reqParam.Dir}}">
                    </div>
                </div>
                <div class="row" style="margin: 10px;">
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
                <div class="row" style="margin: 10px;">
                    <div class="col-md-4 text-right">
                        当前路径：
                    </div>
                    <div class="col-md-6">
                        全部文件<span id="dirPathTextInUploadModal"></span>
                        <input type="hidden" name="parentDir" id="dirPathInUploadModal" readonly value="{{.reqParam.Dir}}">
                    </div>
                </div>
                <div class="row" style="margin: 10px;">
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

<!-- 模态框（Modal） rename -->
<div class="modal fade" id="renameModal" tabindex="-1" role="dialog" aria-labelledby="renameLabel" aria-hidden="true">
    <div class="modal-dialog">
        <div class="modal-content">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal" 
                        aria-hidden="true">×
                </button>
                <h4 class="modal-title" id="renameLabel">
                    重命名
                </h4>
            </div>
            <form id="renameForm" accept-charset="utf-8" role="form" class="form-horizontal" method="POST" action='/my/file/rename' enctype="application/x-www-form-urlencoded">
            <div class="modal-body">
                <div class="row" style="margin: 10px;">
                    <div class="col-md-4 text-right">
                        当前路径：
                    </div>
                    <div class="col-md-6">
                        全部文件<span id="dirPathTextInRenameModal"></span>
                    </div>
                </div>
                <div class="row" style="margin: 10px;">
                    <div class="col-md-4 text-right">
                        原文件名：
                    </div>
                    <div class="col-md-6">
                        <span id="fileNameTextInRenameModal"></span>
                        <input type="hidden" name="id" id="userFileIdInRenameModal" value="">
                    </div>
                </div>
                <div class="row" style="margin: 10px;">
                    <div class="col-md-4 text-right">
                        重命名为：
                    </div>
                    <div class="col-md-6">
                        <input type="text" name="name" id="fileNameInRenameModal" value="">
                    </div>
                </div>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-default" data-dismiss="modal">关闭</button>
                <button type="submit" class="btn btn-primary" id="renameSubmit">提交</button>
            </div>
            </form>
        </div><!-- /.modal-content -->
    </div><!-- /.modal-dialog -->
</div><!-- /.modal -->

<!-- 模态框（Modal） share -->
<div class="modal fade" id="shareModal" tabindex="-1" role="dialog" aria-labelledby="shareLabel" aria-hidden="true">
    <div class="modal-dialog">
        <div class="modal-content">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal" aria-hidden="true">×</button>
                <h4 class="modal-title" id="shareLabel">
                    分享
                </h4>
            </div>
            <form id="uploadForm" accept-charset="utf-8" role="form" class="form-horizontal" method="POST" action='/my/share/add' enctype="application/x-www-form-urlencoded">
            </form>
            <div class="modal-body">
                <div class="row" style="margin: 10px;">
                    <div class="col-md-3 text-right">
                        当前文件：
                    </div>
                    <div class="col-md-7">
                        全部文件<span id="dirPathTextInShareModal"></span><span id="fileNameTextInShareModal"></span>
                        <input type="hidden" name="parentDir" id="dirPathInShareModal" value="{{.reqParam.Dir}}">
                        <input type="hidden" name="fileName" id="fileNameInShareModal" value="">
                        <input type="hidden" name="userFileId" id="userFileIdInShareModal" value="">
                    </div>
                </div>
                <!--

                <div class="row" style="margin: 10px;">
                    <div class="col-md-3 text-right">
                        选择有效期：
                    </div>
                    <div class="col-md-7">
                        <input class="expired-control" type="radio" name="expiredType" id="expiredTypeNone" checked value="0">不分享
                        <input class="expired-control" type="radio" name="expiredType" id="expiredTypeRelative" value="2">分享
                        <br>
                        <input type="text" class="form-control" id="expiredText" name="expiredText" value="" placeholder="点击选择有效期">
                    </div>
                </div>
                -->
                <div class="row" style="margin: 10px;">
                    <div class="col-md-3 text-right">
                        当前状态：
                    </div>
                    <div class="col-md-7">
                        <span id="shareStatusText"></span>
                        <button class="btn btn-info" href="javascript:;" id="btnMakeShare">生成分享链接</button>
                    </div>
                </div>
                <div class="row" style="margin: 10px;">
                    <div class="col-md-3 text-right">
                        分享地址：
                    </div>
                    <div class="col-md-7 input-group">
                        <input class="form-control" type="text" id="shareAddr" readonly value="(尚未分享)">
                        <span class="input-group-btn">
                            <a class="btn btn-info btn-search" href="javascript:;">复制</a>
                        </span>
                    </div>
                </div>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-default" data-dismiss="modal">关闭</button>
                <!--<button type="submit" class="btn btn-primary" id="shareSubmit">提交</button>-->
            </div>
        </div><!-- /.modal-content -->
    </div><!-- /.modal-dialog -->
</div><!-- /.modal -->

<!-- 模态框（Modal） offlineTask -->
<div class="modal fade" id="offlineTaskModal" tabindex="-1" role="dialog" aria-labelledby="offlineTaskLabel" aria-hidden="true">
    <div class="modal-dialog">
        <div class="modal-content">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal" 
                        aria-hidden="true">×
                </button>
                <h4 class="modal-title" id="offlineTaskLabel">
                    新建离线下载任务
                </h4>
            </div>
            <form id="offlineTaskForm" accept-charset="utf-8" role="form" class="form-horizontal" method="POST" action='/my/file/newfolder'>

            <div class="modal-body">
                <div class="row" style="margin: 10px;">
                    <div class="col-md-4 text-right">
                        当前路径：
                    </div>
                    <div class="col-md-6">
                        全部文件<span id="dirPathTextInOfflineTaskModal"></span>
                        <input type="hidden" name="parentDir" id="dirPathInOfflineTaskModal" readonly value="{{.reqParam.Dir}}">
                    </div>
                </div>
                <div class="row" style="margin: 10px;">
                    <div class="col-md-4 text-right">
                        请输入资源地址：
                    </div>
                    <div class="col-md-8">
                        <textarea class="form-control" name="dataurl" id="nameInOfflineTaskModal"></textarea>
                    </div>
                </div>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-default" data-dismiss="modal">关闭</button>
                <button type="submit" class="btn btn-primary" id="offlineTaskSubmit">提交</button>
            </div>
            </form>
        </div><!-- /.modal-content -->
    </div><!-- /.modal-dialog -->
</div><!-- /.modal -->

    
<style type="text/css">
    .glyphicon-folder-close {
        color: #FFCC33; 
    }
    .glyphicon-file {
        color: #68bde6;
    }
</style>

<script type="text/javascript">
Date.prototype.Format = function (fmt) { //author: meizz
  var o = {
    "M+": this.getMonth() + 1, //月份
    "d+": this.getDate(), //日
    "h+": this.getHours(), //小时
    "m+": this.getMinutes(), //分
    "s+": this.getSeconds(), //秒
    "q+": Math.floor((this.getMonth() + 3) / 3), //季度
    "S": this.getMilliseconds() //毫秒
  };
  if (/(y+)/.test(fmt)) fmt = fmt.replace(RegExp.$1, (this.getFullYear() + "").substr(4 - RegExp.$1.length));
  for (var k in o)
  if (new RegExp("(" + k + ")").test(fmt)) fmt = fmt.replace(RegExp.$1, (RegExp.$1.length == 1) ? (o[k]) : (("00" + o[k]).substr(("" + o[k]).length)));
  return fmt;
}

$("#expiredText").datetimepicker({
    format: 'YYYY-MM-DD HH:mm',
});

$(function () {
    // -------- 新建文件夹 ----------
    $('#newFolderModal').on('show.bs.modal', function () {
        $('#dirPathTextInNewFolderModal').html($('#dirPath').val());
        $('#nameInNewFolderModal').focus();// 未生效
    });
    $('#newFolderSubmit').on('click', function(){
        $('#newFolderForm').submit();
        $('#newFolderModal').modal('hide');
    });

    // -------- 上传 ----------
    $('#uploadModal').on('show.bs.modal', function () {
        $('#dirPathTextInUploadModal').html($('#dirPath').val());
    });
    $('#uploadSubmit').on('click', function(){
        $('#uploadForm').submit();
        $('#uploadModal').modal('hide');
    });

    // -------- 重命名 ----------
    $('#renameModal').on('show.bs.modal', function () {
        $('#dirPathTextInRenameModal').html($('#dirPath').val());
        $('#fileTextInRenameModal').html($('#dirPath').val());
    });
    $('#renameSubmit').on('click', function(){
        $('#renameForm').submit();
        $('#renameModal').modal('hide');
    });

    // -------- 分享 ----------
    $('#shareModal').on('show.bs.modal', function () {
        $('#dirPathTextInShareModal').html($('#dirPath').val());
        
        var userFileId = $('#userFileIdInShareModal').val();//$(this).attr('data-id');
        var fileName = $('#fileNameInShareModal').val();//$(this).attr('data-id');
        console.log('userFileId=', userFileId);
        $.ajax({
            url:"/my/share/check/"+userFileId,
            success:function(data, textStatus) {
                console.log(data);
                if (data.result != undefined && data.result.Id != undefined) {
                    // 有分享
                    console.log('the fileid=', data.result.Id);
                    $('#shareStatusText').html('已生成分享链接');
                    $('#shareAddr').val(data.result.OuterPath);
                    $('#btnMakeShare').hide();
                } else {
                    // 无分享
                    $('#shareStatusText').html('尚未分享');
                    $('#btnMakeShare').show();
                }
            }
        });
    });
    // $('#shareSubmit').on('click', function(){
    //     $('#shareForm').submit();
    //     $('#shareModal').modal('hide');
    // });
    $('#btnMakeShare').on('click', function() {
        var userFileId = $('#userFileIdInShareModal').val();
        $.ajax({
            method: 'POST',
            url:"/my/share/add",
            data: {
                'userFileId': userFileId,
            },
            success:function(data, textStatus) {
                console.log(data);
                if (data.result != undefined && data.result.Id != undefined) {
                    // 有分享
                    console.log('the fileid=', data.result.Id);
                    $('#shareStatusText').html('已生成分享链接');
                    $('#shareAddr').val(data.result.OuterPath);
                    $('#btnMakeShare').hide();
                } else {
                    // 无分享
                    alert('生成分享链接失败:'+data.errmsg);
                }
            }
        });

    });

    $('#expiredText').hide();
    // (new Date()).Format("yyyy-M-d h:m:s.S")
    $('#expiredText').val((new Date()).Format("yyyy-MM-dd hh:mm"));
    $('.expired-control').on('click', function(){
        var val = $(this).val();
        console.log('expired-coltrol val=', val);
        if (val == 0) {
            $('#expiredText').hide();
            //$('#expired-text').hide();
        }
        if (val == 2) {
            // 相对有效期
            $('#expiredText').show();
            //$('#expired-text').hide();
        }
    });

    // -------- 离线任务 ----------
    $('#offlineTaskModal').on('show.bs.modal', function () {
        $('#dirPathTextInOfflineTaskModal').html($('#dirPath').val());
        $('#nameInOfflineTaskModal').focus();// 未生效
    });
    $('#offlineTaskSubmit').on('click', function(){
        $('#offlineTaskForm').submit();
        $('#offlineTaskModal').modal('hide');
    });

});

function checkFile(userFileId, fileName) {
    $('#userFileIdInShareModal').val(userFileId);
    $('#fileNameInShareModal').val(fileName);
    $('#fileNameTextInShareModal').html(fileName);

    $('#userFileIdInRenameModal').val(userFileId);
    $('#fileNameInRenameModal').val(fileName);
    $('#fileNameTextInRenameModal').html(fileName);

    //$('#shareModal').show();
}

</script>
    
{{template "common/partfooter.tpl"}}
{{template "common/foot.tpl"}}

{{ end }}
