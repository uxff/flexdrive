{{ define "order/list.tpl" }}

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
            <li><a href="/my/file/list">我的</a></li>
            <li class="active">订单列表</li>
        </ul>
    </div>

    <div class="panel panel-default">
        <div class="panel-heading">查询条件</div>
        <div class="panel-body">
            <form id="formSearch" class="form-horizontal" method="GET" action="/my/order/list">
                <div class="form-group" >
                    <label class="control-label col-sm-1" for="txt_search_name">文件名称</label>
                    <div class="col-sm-2">
                        <input type="text" class="form-control" id="txt_search_name" name="name" value="{{.reqParam.Name}}">
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
                    </div>
                </div>
            </form>
        </div>
    </div>


    <div class="row">
        <table class="table table-striped table-bordered table-hover">
            <thead>
                <tr class="info">
                    <th>订单编号</th>
                    <th>下单时间</th>
                    <th>购买等级</th>
                    <th>获得权益</th>
                    <th>金额(元)</th>
                    <th>状态</th>
                    <th>操作</th>
                </tr>
            </thead>
            <tbody>
                {{range .list}}
                <tr data-id="{{.Id}}" >
                    <td>
                        {{.Id}}
                    </td>
                    <td>{{.Created }}</td>
                    <td>{{.AwardLevelId }}</td>
                    <td>空间增加{{size4Human .AwardSpace }}</td>
                    <td>{{.TotalAmount }}</td>
                    <td>{{orderStatus .Status }}</td>
                    <td>
                        {{if eq .Status 1}}
                            <a href="/" target="_blank">去支付</a>
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
                <button type="submit" class="btn btn-primary" id="shareSubmit">提交</button>
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
                    console.log('the fileid=', data.result.Id);
                }
            }
        });
    });
    $('#shareSubmit').on('click', function(){
        $('#shareForm').submit();
        $('#shareModal').modal('hide');
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

});
    function checkShare(userFileId, fileName) {
        $('#userFileIdInShareModal').val(userFileId);
        $('#fileNameInShareModal').val(fileName);
        $('#fileNameTextInShareModal').val(fileName);
        //$('#shareModal').show();
    }

</script>
    
{{template "common/partfooter.tpl"}}
{{template "common/foot.tpl"}}

{{ end }}
