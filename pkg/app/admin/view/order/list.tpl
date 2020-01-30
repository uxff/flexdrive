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
            <li><a href="/order/list">订单</a></li>
            <li class="active">订单列表</li>
        </ul>
        <div class="panel panel-default">
            <div class="panel-heading">查询条件</div>
            <div class="panel-body">
                <form id="formSearch" class="form-horizontal" method="GET" action="/order/list">
                    <div class="form-group" >
                        <label class="control-label col-sm-1" for="txt_search_userEmail">用户邮箱</label>
                        <div class="col-sm-2">
                            <input type="text" class="form-control" id="txt_search_userEmail" name="userEmail" value="{{.reqParam.UserEmail}}">
                        </div>
                        <label class="control-label col-sm-1" for="txt_search_status">状态</label>
                        <div class="col-sm-2">
                            <select name="status">
                                <option value="0">全部</option>
                                {{range $status, $statusDesc := .orderStatusMap}}
                                <option value="{{$status}}">{{$statusDesc}}</option>
                                {{end}}
                            </select>
                            <input type="text" class="form-control" id="txt_search_status" name="userEmail" value="{{.reqParam.UserEmail}}">
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

        <table class="table table-striped table-bordered table-hover">
            <thead>
                <tr class="info">
                    <th>id</th>
                    <th>会员邮箱</th>
                    <th>会员手机号</th>
                    <th>原等级</th>
                    <th>购买等级</th>
                    <th>获得空间(MB)</th>
                    <th>订单金额(元)</th>
                    <th>下单时间</th>
                    <th>状态</th>
                    <th>操作</th>
                </tr>
            </thead>
            <tbody>
                {{range .list}}
                <tr>
                    <td>{{.Id}}</td>
                    <td>{{.UserId}} ..</td>
                    <td>{{.Phone}} ..</td>
                    <td>{{.OriginLevelId}}</td>
                    <td>{{.AwardLevelId }}</td>
                    <td>{{.AwardSpace }}</td>
                    <td>{{.TotalAmount }}</td>
                    <td>{{.Created }}</td>
                    <td>{{orderStatus .Status}}</td>
                    <td>
                        {{if eq .Status 3}}
                        <a href="/order/refund/{{.Id}}">退款</a>
                        {{end}}
                        {{if eq .Status 4}}
                        <a href="/order/refund/{{.Id}}">重新退款</a>
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
