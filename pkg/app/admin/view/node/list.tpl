{{ define "node/list.tpl" }}

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
            <li><a href="/node/list">节点</a></li>
            <li class="active">节点列表</li>
        </ul>
        <div class="panel panel-default">
            <div class="panel-heading">查询条件</div>
            <div class="panel-body">
                <form id="formSearch" class="form-horizontal" method="GET" action="/node/list">
                    <div class="form-group" >
                        <label class="control-label col-sm-1" for="txt_search_name">名称</label>
                        <div class="col-sm-2">
                            <input type="text" class="form-control" id="txt_search_name" name="name" value="{{.reqParam.Name}}">
                        </div>
                        <label class="control-label col-sm-1" for="txt_search_created">加入时间</label>
                        <div class="col-sm-2">
                            <input type="text" class="form-control" id="txt_search_created_start" name="createStart" value="{{.reqParam.CreateStart}}">
                        </div>
                        <div class="col-sm-2">
                            <input type="text" class="form-control" id="txt_search_created_end" name="createEnd" value="{{.reqParam.CreateEnd}}">
                        </div>
                        <label class="control-label col-sm-1" for="txt_search_last_active">最后注册</label>
                        <div class="col-sm-1">
                            <select class="form-control" id="txt_search_last_active" name="lastActive" value="{{.reqParam.LastActive}}">
                                <option value="-1">不限</option>
                                <option value="60" {{if eq .reqParam.LastActive 60}}selected{{end}}>1分钟内</option>
                                <option value="300" {{if eq .reqParam.LastActive 300}}selected{{end}}>5分钟内</option>
                                <option value="3600" {{if eq .reqParam.LastActive 3600}}selected{{end}}>1小时内</option>
                                <option value="86400" {{if eq .reqParam.LastActive 86400}}selected{{end}}>1天内</option>
                            </select>
                        </div>
                        <div class="col-sm-1" style="text-align:left;">
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
                    <th>名称</th>
                    <th>服务地址</th>
                    <th>已用/总空间</th>
                    <th>文件量</th>
                    <th>加入时间</th>
                    <th>最后心跳时间</th>
                    <th>状态</th>
                    <th>操作</th>
                </tr>
            </thead>
            <tbody>
                {{range $i, $v := .list}}
                <tr>
                    <td>{{$v.Id}}</td>
                    <td>
                        {{if eq $v.Follow $v.NodeName}}
                            <b class="text-success">{{$v.NodeName}} (Master, {{$v.ClusterId}})</b>
                        {{else}}
                            {{$v.NodeName}} (follow: {{$v.Follow}}, {{$v.ClusterId}})
                        {{end}}
                    </td>
                    <td>{{$v.NodeAddr}}</td>
                    <td>{{space4Human $v.UsedSpace }} / {{space4Human $v.TotalSpace }}</td>
                    <td>{{$v.FileCount }}</td>
                    <td>{{$v.Created}}</td>
                    <td style="color:{{timeSmell $v.LastRegistered}}"><b>{{$v.LastRegistered}}</b></td>
                    <td>{{mgrStatus $v.Status}}</td>
                    <td>
                        <button class="btn btn-info setSpace " data-toggle1="modal" data-target1="#setspaceModal" data-id="{{$v.Id}}" data-space="{{$v.TotalSpace}}">设置空间</button>
                        <a href="/node/kick/{{$v.Id}}?ts={{$v.LastRegistered}}">Kick</a>
                        <a href="/node/invite/{{$v.Id}}?ts={{$v.LastRegistered}}">Invite</a>
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

<!-- 模态框（Modal） setspace -->
<div class="modal fade" id="setspaceModal" tabindex="-1" role="dialog" aria-labelledby="setspaceLabel" aria-hidden="true">
    <div class="modal-dialog">
        <div class="modal-content">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal" 
                        aria-hidden="true">×
                </button>
                <h4 class="modal-title" id="setspaceLabel">设置空间</h4>
            </div>
            <form id="setspaceForm" accept-charset="utf-8" role="form" class="form-horizontal" method="POST" action='/node/setspace' enctype="application/x-www-form-urlencoded">
            <div class="modal-body">
                <div class="row" style="margin: 10px;">
                    <div class="col-md-4 text-right">
                        当前节点：
                    </div>
                    <div class="col-md-6">
                        <span id="nodeIdTextInSetspaceModal"></span>
                        <input type="hidden" id="nodeIdInSetspaceModal" name="nodeId" value="">
                    </div>
                </div>
                <div class="row" style="margin: 10px;">
                    <div class="col-md-4 text-right">
                        当前空间：
                    </div>
                    <div class="col-md-6">
                        <span id="totalspaceTextInSetspaceModal"></span>KB
                    </div>
                </div>
                <div class="row" style="margin: 10px;">
                    <div class="col-md-4 text-right">
                        设置为：
                    </div>
                    <div class="col-md-6">
                        <input type="number" name="totalSpace" id="totalspaceInSetspaceModal" value="">KB
                    </div>
                </div>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-default" data-dismiss="modal">关闭</button>
                <button type="submit" class="btn btn-primary" id="setspaceSubmit">提交</button>
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

    // -------- 重size ----------
    $('.setSpace').on('click', function(){
        $('#totalspaceTextInSetspaceModal').html($(this).attr('data-space'));
        $('#totalspaceInSetspaceModal').val($(this).attr('data-space'));
        $('#nodeIdInSetspaceModal').val($(this).attr('data-id'));
        $('#nodeIdTextInSetspaceModal').html($(this).attr('data-id'));
        $('#setspaceModal').modal('show');
    })
    $('#setspaceModal').on('show.bs.modal', function () {
        // $('#dirPathTextInsetspaceModal').html($('#dirPath').val());
        // $('#fileTextInsetspaceModal').html($('#dirPath').val());
    });
    $('#setspaceSubmit').on('click', function(){
        $('#setspaceForm').submit();
        $('#setspaceModal').modal('hide');
    });


</script>
    
{{template "common/partfooter.tpl"}}
{{template "common/foot.tpl"}}

{{ end }}
