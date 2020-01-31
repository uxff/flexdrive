{{ define "share/search.tpl" }}

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
            <li><a href="/share/search">分享</a></li>
            <li class="active">分享搜索</li>
        </ul>

        <form class="form-horizontal" method="GET" action="/share/search">
            <div class="input-group col-md-6" style="margin-top:0px; positon:relative">
                <span class="input-group-addon">分享搜索</span>
                <input type="text" class="form-control" name="name" placeholder="请输入关键字" value="{{.reqParam.Name}}" />
                <span class="input-group-btn">
                    <button class="btn btn-info btn-search" type="submit">🔍搜索</button>
                </span>
            </div>
        </form>
    </div>
    <p></p>
    <div class="row">
        {{range .list}}
        <div class="col-md-8">
            <h3>{{.FileName}}</h3>
            <p>分享人：{{.UserId}} 大小：{{.Size}} 分享时间：</p>
            <p class="text-success">{{.FilePath}}</p>
        </div>
        {{else}}
        <!--应该加载热门关键字-->
        <p></p>
            {{if .reqParam.Name}}
            <div class="col-md-8">暂无数据</div>
            {{else}}
            <div class="col-md-6">
                热门关键字：
                <div class="btn-group btn-group-justified" role="group" aria-label="...">
                    <a href="/share/search?name=java" class="btn btn-default" role="button">java</a>
                    <a href="/share/search?name=golang" class="btn btn-default" role="button">golang</a>
                    <a href="/share/search?name=分布式" class="btn btn-default" role="button">分布式</a>
                    <a href="/share/search?name=python" class="btn btn-default" role="button">python</a>
                    <a href="/share/search?name=C++" class="btn btn-default" role="button">C++</a>
                </div>
            </div>
            {{end}}

        {{end}}
    </div>
    <p></p>
    <div class="row">
        <!--分页-->
        {{template "paginator2.tpl" .}}
    </div>
    <p></p>

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
