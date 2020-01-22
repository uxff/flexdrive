{{ define "common/error.tpl" }}

{{template "common/head.tpl" .}}
{{template "common/partheader.tpl" .}}

<div class="container">
  <div class="row vertical-offset-75">
  	<div class="col-md-6 col-md-offset-3">
  		<div class="panel panel-default" >
		  	<div class="panel-heading text-center">
		    	<h3 class="panel-title"><strong>温馨提示</strong></h3>
		 	  </div>

		  	<div class="panel-body">
          <p style="text-align:center">{{.errMsg}} <a href="javascript:;" onclick="window.history.go(-1)">返回</a></p>
		    </div>

			</div>
		</div>
	</div>
</div>
{{template "common/partfooter.tpl"}}
{{template "common/foot.tpl"}}

{{ end }}