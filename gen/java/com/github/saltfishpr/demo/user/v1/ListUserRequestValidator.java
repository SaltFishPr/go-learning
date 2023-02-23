// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: user/v1/user.proto

package com.github.saltfishpr.demo.user.v1;


/**
* Validates {@code ListUserRequest} protobuf objects.
*/
@SuppressWarnings("all")
public class ListUserRequestValidator implements io.envoyproxy.pgv.ValidatorImpl<com.github.saltfishpr.demo.user.v1.ListUserRequest>{
	public static io.envoyproxy.pgv.ValidatorImpl validatorFor(Class clazz) {
		if (clazz.equals(com.github.saltfishpr.demo.user.v1.ListUserRequest.class)) return new ListUserRequestValidator();
		
		return null;
	}
		
		private final Long OFFSET__GTE = 0L;
	
		
		private final Long LIMIT__LTE = 500L;
		private final Long LIMIT__GTE = 0L;
	
	
	

	public void assertValid(com.github.saltfishpr.demo.user.v1.ListUserRequest proto, io.envoyproxy.pgv.ValidatorIndex index) throws io.envoyproxy.pgv.ValidationException {
	
			io.envoyproxy.pgv.ComparativeValidation.greaterThanOrEqual(".saltfishpr.demo.user.v1.ListUserRequest.offset", proto.getOffset(), OFFSET__GTE, java.util.Comparator.naturalOrder());
	
			io.envoyproxy.pgv.ComparativeValidation.range(".saltfishpr.demo.user.v1.ListUserRequest.limit", proto.getLimit(), null, LIMIT__LTE, null, LIMIT__GTE, java.util.Comparator.naturalOrder());
	
	
	}

}

