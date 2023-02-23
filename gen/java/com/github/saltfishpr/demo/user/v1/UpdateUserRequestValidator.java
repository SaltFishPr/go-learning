// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: user/v1/user.proto

package com.github.saltfishpr.demo.user.v1;


/**
* Validates {@code UpdateUserRequest} protobuf objects.
*/
@SuppressWarnings("all")
public class UpdateUserRequestValidator implements io.envoyproxy.pgv.ValidatorImpl<com.github.saltfishpr.demo.user.v1.UpdateUserRequest>{
	public static io.envoyproxy.pgv.ValidatorImpl validatorFor(Class clazz) {
		if (clazz.equals(com.github.saltfishpr.demo.user.v1.UpdateUserRequest.class)) return new UpdateUserRequestValidator();
		
		if (clazz.equals(com.github.saltfishpr.demo.user.v1.UpdateUserRequest.User.class)) return new UpdateUserRequest_UserValidator();
		return null;
	}
		
	
		
	
	
	

	public void assertValid(com.github.saltfishpr.demo.user.v1.UpdateUserRequest proto, io.envoyproxy.pgv.ValidatorIndex index) throws io.envoyproxy.pgv.ValidationException {
	
		if (proto.hasUser()) {
			io.envoyproxy.pgv.RequiredValidation.required(".saltfishpr.demo.user.v1.UpdateUserRequest.user", proto.getUser());
		} else {
			io.envoyproxy.pgv.RequiredValidation.required(".saltfishpr.demo.user.v1.UpdateUserRequest.user", null);
		};
			// Validate user
			if (proto.hasUser()) index.validatorFor(proto.getUser()).assertValid(proto.getUser());
	
			// Validate mask
			if (proto.hasMask()) index.validatorFor(proto.getMask()).assertValid(proto.getMask());
	
	
	}

/**
	 * Validates {@code UpdateUserRequest_User} protobuf objects.
	 */
	public static class UpdateUserRequest_UserValidator implements io.envoyproxy.pgv.ValidatorImpl<com.github.saltfishpr.demo.user.v1.UpdateUserRequest.User> {
		
	
		
	
		
	
		
	
	
	

	public void assertValid(com.github.saltfishpr.demo.user.v1.UpdateUserRequest.User proto, io.envoyproxy.pgv.ValidatorIndex index) throws io.envoyproxy.pgv.ValidationException {
	// no validation rules for Name

	
			if ( !proto.getUsername().isEmpty() ) {
			io.envoyproxy.pgv.StringValidation.minLength(".saltfishpr.demo.user.v1.UpdateUserRequest.User.username", proto.getUsername(), 3);
			io.envoyproxy.pgv.StringValidation.maxLength(".saltfishpr.demo.user.v1.UpdateUserRequest.User.username", proto.getUsername(), 32);
			}
	
			if ( !proto.getPassword().isEmpty() ) {
			io.envoyproxy.pgv.StringValidation.minLength(".saltfishpr.demo.user.v1.UpdateUserRequest.User.password", proto.getPassword(), 6);
			io.envoyproxy.pgv.StringValidation.maxLength(".saltfishpr.demo.user.v1.UpdateUserRequest.User.password", proto.getPassword(), 32);
			}
	
			if ( !proto.getEmail().isEmpty() ) {
			io.envoyproxy.pgv.StringValidation.email(".saltfishpr.demo.user.v1.UpdateUserRequest.User.email", proto.getEmail());
			}
	
	
	}
}
}

