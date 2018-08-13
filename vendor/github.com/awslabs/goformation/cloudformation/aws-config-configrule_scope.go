package cloudformation

// AWSConfigConfigRule_Scope AWS CloudFormation Resource (AWS::Config::ConfigRule.Scope)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-config-configrule-scope.html
type AWSConfigConfigRule_Scope struct {

	// ComplianceResourceId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-config-configrule-scope.html#cfn-config-configrule-scope-complianceresourceid
	ComplianceResourceId Value `json:"ComplianceResourceId,omitempty"`

	// ComplianceResourceTypes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-config-configrule-scope.html#cfn-config-configrule-scope-complianceresourcetypes
	ComplianceResourceTypes []Value `json:"ComplianceResourceTypes,omitempty"`

	// TagKey AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-config-configrule-scope.html#cfn-config-configrule-scope-tagkey
	TagKey Value `json:"TagKey,omitempty"`

	// TagValue AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-config-configrule-scope.html#cfn-config-configrule-scope-tagvalue
	TagValue Value `json:"TagValue,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSConfigConfigRule_Scope) AWSCloudFormationType() string {
	return "AWS::Config::ConfigRule.Scope"
}
