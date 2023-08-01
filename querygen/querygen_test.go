package querygen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenApplicationQuery(t *testing.T) {
	fields := []string{
		"objectRef.apiGroup",
		"objectRef.apiVersion",
		"objectRef.resource",
		"verb",
		"auditID",
		"objectRef.name",
		"requestReceivedTimestamp",
		"impersonatedUser.username",
		"user.username",
	}
	json_properties := map[string]string{
		"apiGroup":   "objectRef.apiGroup",
		"apiVersion": "objectRef.apiVersion",
		"kind":       "objectRef.resource",
		"name":       "objectRef.name",
	}
	track_fields := TrackFieldSpec{
		with_userid:     true,
		with_ev_verb:    true,
		with_ev_subject: true,
	}
	expected := `search ` +
		`index="some_index" ` +
		`log_type=audit ` +
		`verb=create ` +
		`"responseStatus.code" IN (200, 201) ` +
		`"objectRef.apiGroup"="appstudio.redhat.com" ` +
		`"objectRef.resource"="applications" ` +
		`("impersonatedUser.username"="*" OR (user.username="*" AND NOT user.username="system:*")) ` +
		`(verb!=create OR "responseObject.metadata.resourceVersion"="*")` +
		`|` + GenDedupEval(fields) +
		`|` + GenTrackFields(track_fields, json_properties)
	out := GenApplicationQuery("some_index")
	assert.Equal(t, expected , out)
}

func TestGenPropertiesJSONExpr(t *testing.T) {
	type args struct {
		properties_map map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Application query properties",
			args: args{properties_map: map[string]string{
				"apiGroup":   "objectRef.apiGroup",
				"apiVersion": "objectRef.apiVersion",
				"kind":       "objectRef.resource",
				"name":       "objectRef.name",
			}},
			want: `json_object(` +
				`"apiGroup",'objectRef.apiGroup',` +
				`"apiVersion",'objectRef.apiVersion',` +
				`"kind",'objectRef.resource',` +
				`"name",'objectRef.name'` +
				`)`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenPropertiesJSONExpr(tt.args.properties_map)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGenTrackFields(t *testing.T) {
	json_properties := map[string]string{
		"apiGroup":   "objectRef.apiGroup",
		"apiVersion": "objectRef.apiVersion",
		"kind":       "objectRef.resource",
		"name":       "objectRef.name",
	}
	type args struct {
		spec           TrackFieldSpec
		properties_map map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "with userid, verb, subject & app props",
			args: args{
				spec: TrackFieldSpec{
					with_userid:     true,
					with_ev_verb:    true,
					with_ev_subject: true,
				},
				properties_map: json_properties,
			},
			want: `eval ` +
				`event_subject='objectRef.resource',` +
				`event_verb='verb',` +
				`messageId='auditID',` +
				`timestamp='requestReceivedTimestamp',` +
				`type="track",` +
				`userId=` + UserIdExpr + `,` +
				`properties=` + GenPropertiesJSONExpr(json_properties) +
				`|` + IncludeFieldsCmd + `|` + ExcludeFieldsCmd,
		},
		{
			name: "with namespace, verb, subject & app props",
			args: args{
				spec: TrackFieldSpec{
					with_namespace:  true,
					with_ev_verb:    true,
					with_ev_subject: true,
				},
				properties_map: json_properties,
			},
			want: `eval ` +
				`event_subject='objectRef.resource',` +
				`event_verb='verb',` +
				`messageId='auditID',` +
				`namespace='objectRef.namespace',` +
				`timestamp='requestReceivedTimestamp',` +
				`type="track",` +
				`properties=` + GenPropertiesJSONExpr(json_properties) +
				`|` + IncludeFieldsCmd + `|` + ExcludeFieldsCmd,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenTrackFields(tt.args.spec, tt.args.properties_map)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGenPipelineRunQuery(t *testing.T) {
	fields := []string{
		"objectRef.apiGroup",
		"objectRef.apiVersion",
		"responseObject.metadata.labels.appstudio.openshift.io/application",
		"responseObject.metadata.labels.appstudio.openshift.io/component",
		"objectRef.resource",
		"verb",
		"auditID",
		"objectRef.namespace",
		"requestReceivedTimestamp",
	}
	json_properties := map[string]string{
		"apiGroup":    "objectRef.apiGroup",
		"apiVersion":  "objectRef.apiVersion",
		"kind":        "objectRef.resource",
		"application": "responseObject.metadata.labels.appstudio.openshift.io/application",
		"component":   "responseObject.metadata.labels.appstudio.openshift.io/component",
	}
	track_fields := TrackFieldSpec{
		with_namespace:  true,
		with_ev_verb:    true,
		with_ev_subject: true,
	}
	expected := `search ` +
		`index="some_index" ` +
		`log_type=audit ` +
		`verb=create ` +
		`"responseStatus.code" IN (200, 201) ` +
		`"objectRef.apiGroup"="tekton.dev" ` +
		`"objectRef.resource"="pipelineruns" ` +
		`"responseObject.metadata.labels.pipelines.appstudio.openshift.io/type"=build` +
		`"responseObject.metadata.resourceVersion"="*"` +
		`|` + GenDedupEval(fields) +
		`|` + GenTrackFields(track_fields, json_properties)
	out := GenPipelineRunQuery("some_index")
	assert.Equal(t, expected, out)
}
