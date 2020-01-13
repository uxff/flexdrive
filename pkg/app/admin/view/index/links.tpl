{{append . "HeadStyles" "/static/css/custom.css"}}
{{append . "HeadScripts" "/static/js/custom.js"}}


<div class="container">
    <div class="row" style="display: none">

    {{template "alert.tpl" .}}
        <div class="btn-group btn-group-justified" role="group" aria-label="...">
            <a href="javascript:;" class="btn btn-default" role="button">广告位招商 A</a>
            <a href="javascript:;" class="btn btn-default" role="button">广告位招商 B</a>
            <a href="javascript:;" class="btn btn-default" role="button">广告位招商 C</a>
        </div>
        <p></p>
    </div>

    <div class="row">

        {{range $gi, $lister := .thelinks}}

            <div class="panel panel-success" style="{{if eq $lister.Hide 1}}display:none{{end}}">
                <div class="panel-heading">
                    <h3 class="panel-title">{{$lister.Name}}</h3>
                </div>
                <div class="panel-body">

                    <div class="row clearfix">
                    {{range $k, $site := $lister.Links}}
                        <div class="col-md-3">
                            <a href="{{$site.Url}}" target="_blank">{{$site.Name}}</a>
                        </div>
                    {{end}}
                    </div>
                </div>
            </div>

        {{end}}
    </div>
</div>
