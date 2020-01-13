{{append . "HeadStyles" "/static/css/custom.css"}}
{{append . "HeadScripts" "/static/js/custom.js"}}


<div class="container">
        <div class="row">
            <div class="panel panel-default">
                <div class="panel-heading">
                    <div class="text-muted " style="display:inline">您的位置：</div>
                    <a href="/picset">图集首页</a>/{{if ne .fullParentName "."}}<a href="{{.parentLink}}">{{.fullParentName}}</a>/{{end}}{{if ne .curDirName "."}}{{.curDirName}}{{end}}
                </div>
                <div class="bootstrap-admin-panel-content span3 arch-warp" style="padding:10px">
                    <div class="row">
                    {{range $k, $aname := .thedirnames}}
                        <div class="col-sm-6 col-md-4">
                            <div class="thumbnail">
                                <a href="{{$aname.Url}}">
                                    <img src="{{$aname.Thumb}}" alt="{{$aname.Name}}" style="height:100%;">
                                </a>
                            </div>
                            <div class="thumbnail-title">
                                <p style="text-align:center;">{{$aname.Name}}</p>
                            </div>
                        </div>
                    {{end}}
                    </div>
                </div>
            </div>
        </div>
{{template "layouts/paginator.html" .}}

</div>


{{/*<script src="/static/js/bootstrap.min.js"></script>*/}}
{{/*<script type="text/javascript" src="/static/js/twitter-bootstrap-hover-dropdown.min.js"></script>*/}}
{{/*<script type="text/javascript" src="/static/js/bootstrap-admin-theme-change-size.js"></script>*/}}
