{{ define "manager/list.tpl" }}

{{template "common/head.tpl" .}}
{{template "common/partheader.tpl" .}}


<div class="container">
    <div class="row" style="display: none">

        <div class="btn-group btn-group-justified" role="group" aria-label="...">
            <a href="javascript:;" class="btn btn-default" role="button">广告位招商 A</a>
            <a href="javascript:;" class="btn btn-default" role="button">广告位招商 B</a>
            <a href="javascript:;" class="btn btn-default" role="button">广告位招商 C</a>
        </div>
        <p></p>
    </div>

    <div class="row">
            <div class="panel panel-success">
                <div class="panel-heading">
                    <h3 class="panel-title">This is a demo site.</h3>
                </div>

                <div class="panel-body">

                    <div class="row clearfix">
                        <div class="col-md-3">
                            Here is the picset demo preview:
                        </div>
                    </div>

                    <img src="https://raw.githubusercontent.com/uxff/taniago/master/20181127073913.png">
                </div>
            </div>
    </div>

    <div class="row">
        <a href="https://github.com/uxff/taniago">view source</a>
    </div>
</div>

{{template "common/partfooter.tpl"}}
{{template "common/foot.tpl"}}

{{ end }}
