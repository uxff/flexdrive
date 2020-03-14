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

    <div class="row">
        <div class="panel panel-default">
            <div class="panel-heading">查询条件</div>
            <div class="panel-body">
                <form id="formSearch" class="form-horizontal" method="GET" action="/my/order/list">
                    <div class="form-group" >
                        <label class="control-label col-sm-1" for="txt_search_name">名称</label>
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
    </div>



    <div class="row">
        <table class="table table-striped table-bordered table-hover">
            <thead>
                <tr class="info">
                    <th>订单编号</th>
                    <th>下单时间</th>
                    <th>购买升级包</th>
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
                    <td>{{.LevelName}}({{.AwardLevelId }})</td>
                    <td>空间增加 {{space4Human .AwardSpace }}</td>
                    <td>{{.TotalAmount }}</td>
                    <td>{{orderStatus .Status }}</td>
                    <td>
                        <a href="/my/order/detail/{{.Id}}" target="_blank">详情</a>
                        {{if eq .Status 1}}
                            <a href="/my/order/mockpay/{{.Id}}" target="_blank">去支付</a>
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

$("#txt_search_created_start").datetimepicker({
    format: 'YYYY-MM-DD HH:mm'
});

$("#txt_search_created_end").datetimepicker({
    format: 'YYYY-MM-DD HH:mm'
});

$(function () {

    $('#expiredText').val((new Date()).Format("yyyy-MM-dd hh:mm"));

});

</script>
    
{{template "common/partfooter.tpl"}}
{{template "common/foot.tpl"}}

{{ end }}
