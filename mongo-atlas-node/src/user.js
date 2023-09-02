const Joi = require('joi');

const userSchema = Joi.object({
  first: Joi.string().required().messages({
    'string.empty': 'First name is required',
    'any.required': 'First name is required'
  }),
  last: Joi.string().required().messages({
    'string.empty': 'Last name is required',
    'any.required': 'Last name is required'
  }),
  email: Joi.string().email().required().messages({
    'string.email': 'Email must be a valid email address',
    'string.empty': 'Email is required',
    'any.required': 'Email is required'
  }),
  status: Joi.string().required().messages({
    'string.empty': 'Status is required',
    'any.required': 'Status is required'
  }),
  city: Joi.string().optional(),
  country: Joi.string().optional(),
  age: Joi.number().integer().min(0).required().messages({
    'number.base': 'Age must be a number',
    'number.integer': 'Age must be an integer',
    'number.min': 'Age must be at least 0',
    'any.required': 'Age is required'
  })
});

class User {
    constructor(user) {
      this._validate(user);
      this.status = 'active'
      this.created_at = new Date();
      this.update_at = new Date();

    }
  
    _validate(user) {
      const { error, value } = userSchema.validate(user)
      if (error) {
        throw new Error(error);
      }
      
      this.first = value.first;
      this.last = value.last;
      this.email = value.email;
      this.age  = value.age;

      return 
    }
  }
  
  module.exports = { User, userSchema };