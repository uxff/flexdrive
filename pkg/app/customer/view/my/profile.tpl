{{ define "my/profile.tpl" }}

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
        <div class="col-md-3" >
            会员账号
        </div>
        <div class="col-md-3" >
            <div class="well well-sm">
                {{.LoginInfo.UserEnt.Email}}
            </div>
        </div>
        <div class="col-md-2"></div>
    </div>
    <div class="row">
        <div class="col-md-2"></div>
        <div class="col-md-3" >
            账号id
        </div>
        <div class="col-md-3" >
            <div class="well well-sm">
                {{.LoginInfo.UserEnt.Id}}
            </div>
        </div>
        <div class="col-md-2"></div>
    </div>
    <div class="row hide">
        <div class="col-md-2"></div>
        <div class="col-md-3" >
            会员升级包
        </div>
        <div class="col-md-3" >
            <div class="well well-sm">
            {{.userLevel.Name}} 
            <a href="/my/order/create" class="btn btn-info  btn-sm" type="button">购买升级包</a>
            </div>
        </div>
        <div class="col-md-2"></div>
    </div>
    <div class="row">
        <div class="col-md-2"></div>
        <div class="col-md-3" >
            总空间
        </div>
        <div class="col-md-3" >
            <div class="well well-sm">
                {{space4Human .LoginInfo.UserEnt.QuotaSpace}}
            </div>
        </div>
        <div class="col-md-2"><a href="/my/order/create" class="btn btn-info btn-sm" style="text-align: right;" >购买升级</a></div>
    </div>
    <div class="row">
        <div class="col-md-2"></div>
        <div class="col-md-3" >
            已用空间
        </div>
        <div class="col-md-3" >
            <div class="well well-sm">
                {{space4Human .LoginInfo.UserEnt.UsedSpace}}
                
            </div>
        </div>
        <div class="col-md-2"></div>
    </div>
    <div class="row" style="margin-bottom: 4px;">
        <div class="col-md-2"></div>
        <div class="col-md-3" >
            使用比例
        </div>
        <div class="col-md-3" >
            <div class="progress " style="width:100%; float: left; height: 23px; margin-bottom: 10px; background-color: #dff0d8;">
                <div class="progress-bar progress-bar-success" role="progressbar"
                        aria-valuenow="60" aria-valuemin="0" aria-valuemax="100"
                        style="width: {{spaceRate .LoginInfo.UserEnt.UsedSpace .LoginInfo.UserEnt.QuotaSpace}}%;">
                    <span class="sr-only">{{spaceRate .LoginInfo.UserEnt.UsedSpace .LoginInfo.UserEnt.QuotaSpace}}% 已用</span>
                </div>
                {{spaceRate .LoginInfo.UserEnt.UsedSpace .LoginInfo.UserEnt.QuotaSpace}}%
            </div>
        </div>
        <div class="col-md-2"></div>
    </div>
    <div class="row">
        <div class="col-md-2"></div>
        <div class="col-md-3" >
            总文件数
        </div>
        <div class="col-md-3" >
            <div class="well well-sm">
                {{.LoginInfo.UserEnt.FileCount}}
            </div>
        </div>
        <div class="col-md-2"><a href="/my/file/list" class="btn btn-info btn-sm" style="text-align: right;" >文件列表</a></div>
    </div>
    <div class="row">
        <div class="col-md-2"></div>
        <div class="col-md-3" >
            分享数量
        </div>
        <div class="col-md-3" >
            <div class="well well-sm">
                2
            </div>
        </div>
        <div class="col-md-2"><a href="/my/share/list" class="btn btn-info btn-sm" style="text-align: right;" >分享列表</a></div>
    </div>
    <div class="row">
        <div class="col-md-2"></div>
        <div class="col-md-3" >
            累计充值
        </div>
        <div class="col-md-3" >
            <div class="well well-sm">
                {{amount4Human .LoginInfo.UserEnt.TotalCharge}} 元
                
            </div>
        </div>
        <div class="col-md-2"><a href="/my/order/list" class="btn btn-info btn-sm" style="text-align: right;" >订单列表</a></div>
    </div>
    <div class="row">
        <div class="col-md-2"></div>
        <div class="col-md-3" >
            最后登录时间
        </div>
        <div class="col-md-3" >
            <div class="well well-sm">
                {{.LoginInfo.UserEnt.LastLoginAt}}
            </div>
        </div>
        <div class="col-md-2"></div>
    </div>
    <div class="row">
        <div class="col-md-2"></div>
        <div class="col-md-3" >
            最后登录ip
        </div>
        <div class="col-md-3" >
            <div class="well well-sm">
                {{.LoginInfo.UserEnt.LastLoginIp}}
            </div>
        </div>
        <div class="col-md-2"></div>
    </div>
    


</div>
 

    
{{template "common/partfooter.tpl"}}
{{template "common/foot.tpl"}}

{{ end }}
