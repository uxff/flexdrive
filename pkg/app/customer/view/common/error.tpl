{{ define "common/error.tpl" }}

{{template "common/head.tpl" .}}
{{template "common/partheader.tpl" .}}

<div class="container">
  <div class="row" style="margin-top: 80px;margin-bottom: 80px;">
  	<div class="col-md-6 col-md-offset-3">
  		<div class="panel panel-danger">
		  	<div class="panel-heading text-center">
		    	<h3 class="panel-title"><strong>温馨提示</strong></h3>
			</div>

		  	<div class="panel-body" style="padding-top: 30px;padding-bottom: 30px;">
				<p style="text-align:center">
					{{.errMsg}}&nbsp;&nbsp;
					<a href="javascript:;" onclick="window.history.go(-1)">返回</a>&nbsp;&nbsp;
					{{if not .IsLogin}}
					<a href="/login" >登录</a>
					{{end}}
				</p>
		    </div>

			</div>
		</div>
	</div>
</div>
{{template "common/partfooter.tpl"}}
{{template "common/foot.tpl"}}

{{ end }}