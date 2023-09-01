class User {
    constructor(user) {
      this._validate(user);
      this.status = 'active'
      this.created_at = new Date();
      this.update_at = new Date();

    }
  
    _validate(user) {
      const { first, last, email, age } = user;
  
      if (!first || !last || !email || !age ) {
        throw new Error('first, last, email and age are required');
      }
  
      if (typeof first !== 'string' || typeof last !== 'string' || typeof email !== 'string' || typeof age !== 'number' ) {
        throw new Error('first, last and email must be strings and age must be a number');
      }
  
      if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
        throw new Error('Invalid email format');
      }
      this.first = first;
      this.last = last;
      this.email = email;
      this.age  = age;

      return 
    }
  }
  
  module.exports = User;