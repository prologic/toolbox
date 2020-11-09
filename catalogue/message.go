package catalogue

// Code generated by dev catalogue command DO NOT EDIT

import (
	domaindropboxapidbx_auth "github.com/watermint/toolbox/domain/dropbox/api/dbx_auth"
	domaindropboxapidbx_auth_attr "github.com/watermint/toolbox/domain/dropbox/api/dbx_auth_attr"
	domaindropboxapidbx_conn_impl "github.com/watermint/toolbox/domain/dropbox/api/dbx_conn_impl"
	domaindropboxapidbx_list_impl "github.com/watermint/toolbox/domain/dropbox/api/dbx_list_impl"
	domaindropboxmodelmo_file_filter "github.com/watermint/toolbox/domain/dropbox/model/mo_file_filter"
	domaindropboxmodelmo_sharedfolder_member "github.com/watermint/toolbox/domain/dropbox/model/mo_sharedfolder_member"
	domaindropboxusecaseuc_compare_local "github.com/watermint/toolbox/domain/dropbox/usecase/uc_compare_local"
	domaindropboxusecaseuc_compare_paths "github.com/watermint/toolbox/domain/dropbox/usecase/uc_compare_paths"
	domaindropboxusecaseuc_file_merge "github.com/watermint/toolbox/domain/dropbox/usecase/uc_file_merge"
	domaindropboxusecaseuc_file_relocation "github.com/watermint/toolbox/domain/dropbox/usecase/uc_file_relocation"
	domaingoogleservicesv_label "github.com/watermint/toolbox/domain/google/service/sv_label"
	domaingoogleservicesv_message "github.com/watermint/toolbox/domain/google/service/sv_message"
	essentialslogesl_rotate "github.com/watermint/toolbox/essentials/log/esl_rotate"
	essentialsmodelmo_filter "github.com/watermint/toolbox/essentials/model/mo_filter"
	essentialsnetworknw_diag "github.com/watermint/toolbox/essentials/network/nw_diag"
	infraapiapi_auth_impl "github.com/watermint/toolbox/infra/api/api_auth_impl"
	infraapiapi_callback "github.com/watermint/toolbox/infra/api/api_callback"
	infracontrolapp_error "github.com/watermint/toolbox/infra/control/app_error"
	infracontrolapp_job_impl "github.com/watermint/toolbox/infra/control/app_job_impl"
	infradocdc_options "github.com/watermint/toolbox/infra/doc/dc_options"
	infrafeedfd_file_impl "github.com/watermint/toolbox/infra/feed/fd_file_impl"
	infrareciperc_group "github.com/watermint/toolbox/infra/recipe/rc_group"
	infrareciperc_group_impl "github.com/watermint/toolbox/infra/recipe/rc_group_impl"
	infrareciperc_spec "github.com/watermint/toolbox/infra/recipe/rc_spec"
	infrareciperc_value "github.com/watermint/toolbox/infra/recipe/rc_value"
	infrareportrp_model_impl "github.com/watermint/toolbox/infra/report/rp_model_impl"
	infrareportrp_writer_impl "github.com/watermint/toolbox/infra/report/rp_writer_impl"
	infrauiapp_ui "github.com/watermint/toolbox/infra/ui/app_ui"
	ingredientfile "github.com/watermint/toolbox/ingredient/file"
	recipedevdiag "github.com/watermint/toolbox/recipe/dev/diag"
	recipefile "github.com/watermint/toolbox/recipe/file"
	recipefiledispatch "github.com/watermint/toolbox/recipe/file/dispatch"
	recipefileimportbatch "github.com/watermint/toolbox/recipe/file/import/batch"
	recipegroupmember "github.com/watermint/toolbox/recipe/group/member"
	recipemember "github.com/watermint/toolbox/recipe/member"
	recipememberquota "github.com/watermint/toolbox/recipe/member/quota"
	recipememberupdate "github.com/watermint/toolbox/recipe/member/update"
	recipeservicesgithubreleaseasset "github.com/watermint/toolbox/recipe/services/github/release/asset"
	recipesharedfoldermember "github.com/watermint/toolbox/recipe/sharedfolder/member"
	recipeteamactivity "github.com/watermint/toolbox/recipe/team/activity"
	recipeteamactivitybatch "github.com/watermint/toolbox/recipe/team/activity/batch"
	recipeteamdevice "github.com/watermint/toolbox/recipe/team/device"
	recipeteamfilerequest "github.com/watermint/toolbox/recipe/team/filerequest"
	recipeteamnamespacemember "github.com/watermint/toolbox/recipe/team/namespace/member"
	recipeteamsharedlink "github.com/watermint/toolbox/recipe/team/sharedlink"
	recipeteamsharedlinkupdate "github.com/watermint/toolbox/recipe/team/sharedlink/update"
)

func AutoDetectedMessageObjects() []interface{} {
	return []interface{}{
		&domaindropboxapidbx_auth.MsgGenerated{},
		&domaindropboxapidbx_auth_attr.MsgAttr{},
		&domaindropboxapidbx_conn_impl.MsgConnect{},
		&domaindropboxapidbx_list_impl.MsgList{},
		&domaindropboxmodelmo_file_filter.MsgFileFilterOpt{},
		&domaindropboxmodelmo_sharedfolder_member.MsgExternalOpt{},
		&domaindropboxmodelmo_sharedfolder_member.MsgInternalOpt{},
		&domaindropboxusecaseuc_compare_local.MsgCompare{},
		&domaindropboxusecaseuc_compare_paths.MsgCompare{},
		&domaindropboxusecaseuc_file_merge.MsgMerge{},
		&domaindropboxusecaseuc_file_relocation.MsgRelocation{},
		&domaingoogleservicesv_label.MsgFindLabel{},
		&domaingoogleservicesv_message.MsgProgress{},
		&essentialslogesl_rotate.MsgOut{},
		&essentialslogesl_rotate.MsgPurge{},
		&essentialslogesl_rotate.MsgRotate{},
		&essentialsmodelmo_filter.MsgFilter{},
		&essentialsnetworknw_diag.MsgNetwork{},
		&infraapiapi_auth_impl.MsgApiAuth{},
		&infraapiapi_callback.MsgCallback{},
		&infracontrolapp_error.MsgErrorReport{},
		&infracontrolapp_job_impl.MsgLauncher{},
		&infradocdc_options.MsgDoc{},
		&infrafeedfd_file_impl.MsgRowFeed{},
		&infrareciperc_group.MsgHeader{},
		&infrareciperc_group_impl.MsgGroup{},
		&infrareciperc_spec.MsgSelfContained{},
		&infrareciperc_value.MsgRepository{},
		&infrareciperc_value.MsgValFdFileRowFeed{},
		&infrareportrp_model_impl.MsgColumnSpec{},
		&infrareportrp_model_impl.MsgTransactionReport{},
		&infrareportrp_writer_impl.MsgSortedWriter{},
		&infrareportrp_writer_impl.MsgXlsxWriter{},
		&infrauiapp_ui.MsgConsole{},
		&infrauiapp_ui.MsgProgress{},
		&ingredientfile.MsgUpload{},
		&recipedevdiag.MsgLoader{},
		&recipefile.MsgRestore{},
		&recipefiledispatch.MsgLocal{},
		&recipefileimportbatch.MsgUrl{},
		&recipegroupmember.MsgList{},
		&recipemember.MsgInvite{},
		&recipememberquota.MsgList{},
		&recipememberquota.MsgUpdate{},
		&recipememberquota.MsgUsage{},
		&recipememberupdate.MsgEmail{},
		&recipeservicesgithubreleaseasset.MsgUp{},
		&recipesharedfoldermember.MsgList{},
		&recipeteamactivity.MsgUser{},
		&recipeteamactivitybatch.MsgUser{},
		&recipeteamdevice.MsgUnlink{},
		&recipeteamfilerequest.MsgList{},
		&recipeteamnamespacemember.MsgList{},
		&recipeteamsharedlink.MsgList{},
		&recipeteamsharedlinkupdate.MsgExpiry{},
	}
}
