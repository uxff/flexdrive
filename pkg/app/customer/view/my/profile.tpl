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
            <li class="active">我的主页</li>
        </ul>
    </div>

    <div class="row">
        <div class="col-md-2"></div>
        <div class="col-md-4" >
            会员账号
        </div>
        <div class="col-md-4" >
            {{.LoginInfo.UserEnt.Email}}
        </div>
        <div class="col-md-2"></div>
    </div>
    <div class="row">
        <div class="col-md-2"></div>
        <div class="col-md-4" >
            会员等级
        </div>
        <div class="col-md-4" >
            {{.LoginInfo.UserEnt.Email}}
        </div>
        <div class="col-md-2"></div>
    </div>
    <div class="row">
        <div class="col-md-2"></div>
        <div class="col-md-4" >
            会员等级
        </div>
        <div class="col-md-4" >
            {{.userLevel.Name}} <button class="btn btn-info " type="button">升级</button>
        </div>
        <div class="col-md-2"></div>
    </div>
    <div class="row">
        <div class="col-md-2"></div>
        <div class="col-md-4" >
            空间
        </div>
        <div class="col-md-4" >
            已用 {{space4Human .LoginInfo.UserEnt.UsedSpace}} / 总共 {{space4Human .LoginInfo.UserEnt.QuotaSpace}}
            [<a href="/my/file/list" style="text-align: right;" >文件列表</a>]
            <div class="progress " style="width:100%; float: left; height: 6px; margin-bottom: 10px; background-color: #dff0d8;">
                <div class="progress-bar progress-bar-success" role="progressbar"
                        aria-valuenow="60" aria-valuemin="0" aria-valuemax="100"
                        style="width: {{spaceRate .LoginInfo.UserEnt.UsedSpace .LoginInfo.UserEnt.QuotaSpace}}%;">
                    <span class="sr-only">{{spaceRate .LoginInfo.UserEnt.UsedSpace .LoginInfo.UserEnt.QuotaSpace}}% 已用</span>
                </div>
            </div>
        </div>
        <div class="col-md-2"></div>
    </div>
    <div class="row" style="margin-bottom: 4px;">
        <div class="col-md-2"></div>
        <div class="col-md-8" >
            当前等级：{{.userLevel.Name}} [<a href="/" >升级</a>]&nbsp;&nbsp;
            当前空间：已用 {{space4Human .LoginInfo.UserEnt.UsedSpace}} / 总共 {{space4Human .LoginInfo.UserEnt.QuotaSpace}}
            [<a href="/" style="text-align: right;" >扩容</a>]
            <div class="progress " style="width:100%; float: left; height: 6px; margin-bottom: 10px; background-color: #dff0d8;">
                <div class="progress-bar progress-bar-success" role="progressbar"
                        aria-valuenow="60" aria-valuemin="0" aria-valuemax="100"
                        style="width: {{spaceRate .LoginInfo.UserEnt.UsedSpace .LoginInfo.UserEnt.QuotaSpace}}%;">
                    <span class="sr-only">{{spaceRate .LoginInfo.UserEnt.UsedSpace .LoginInfo.UserEnt.QuotaSpace}}% 已用</span>
                </div>
            </div>
        </div>
        <div class="col-md-2"></div>
    </div>
    <div class="row" style="margin-bottom: 4px;">
        <div class="col-md-2"></div>
        <div class="col-md-8" >
            当前等级：{{.userLevel.Name}} [<a href="/" >升级</a>]&nbsp;&nbsp;
            当前空间：已用 {{space4Human .LoginInfo.UserEnt.UsedSpace}} / 总共 {{space4Human .LoginInfo.UserEnt.QuotaSpace}}
            [<a href="/" style="text-align: right;" >扩容</a>]
            <div class="progress " style="width:100%; float: left; height: 6px; margin-bottom: 10px; background-color: #dff0d8;">
                <div class="progress-bar progress-bar-success" role="progressbar"
                        aria-valuenow="60" aria-valuemin="0" aria-valuemax="100"
                        style="width: {{spaceRate .LoginInfo.UserEnt.UsedSpace .LoginInfo.UserEnt.QuotaSpace}}%;">
                    <span class="sr-only">{{spaceRate .LoginInfo.UserEnt.UsedSpace .LoginInfo.UserEnt.QuotaSpace}}% 已用</span>
                </div>
            </div>
        </div>
        <div class="col-md-2"></div>
    </div>
    <div class="row" style="margin-bottom: 4px;">
        <div class="col-md-2"></div>
        <div class="col-md-8" >
            当前等级：{{.userLevel.Name}} [<a href="/" >升级</a>]&nbsp;&nbsp;
            当前空间：已用 {{space4Human .LoginInfo.UserEnt.UsedSpace}} / 总共 {{space4Human .LoginInfo.UserEnt.QuotaSpace}}
            [<a href="/" style="text-align: right;" >扩容</a>]
            <div class="progress " style="width:100%; float: left; height: 6px; margin-bottom: 10px; background-color: #dff0d8;">
                <div class="progress-bar progress-bar-success" role="progressbar"
                        aria-valuenow="60" aria-valuemin="0" aria-valuemax="100"
                        style="width: {{spaceRate .LoginInfo.UserEnt.UsedSpace .LoginInfo.UserEnt.QuotaSpace}}%;">
                    <span class="sr-only">{{spaceRate .LoginInfo.UserEnt.UsedSpace .LoginInfo.UserEnt.QuotaSpace}}% 已用</span>
                </div>
            </div>
        </div>
        <div class="col-md-2"></div>
    </div>



</div>
 
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
