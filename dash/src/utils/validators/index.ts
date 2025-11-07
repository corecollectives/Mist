
import { APP_CONSTANTS } from '../../constants';

export interface ValidationResult {
  isValid: boolean;
  message?: string;
}

export function validateEmail(email: string): ValidationResult {
  if (!email) {
    return { isValid: false, message: 'Email is required' };
  }
  
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (!emailRegex.test(email)) {
    return { isValid: false, message: 'Please enter a valid email address' };
  }
  
  return { isValid: true };
}

export function validatePassword(password: string): ValidationResult {
  if (!password) {
    return { isValid: false, message: 'Password is required' };
  }
  
  if (password.length < APP_CONSTANTS.VALIDATION.MIN_PASSWORD_LENGTH) {
    return { 
      isValid: false, 
      message: `Password must be at least ${APP_CONSTANTS.VALIDATION.MIN_PASSWORD_LENGTH} characters long` 
    };
  }
  
  const hasUppercase = /[A-Z]/.test(password);
  const hasLowercase = /[a-z]/.test(password);
  const hasNumber = /\d/.test(password);
  
  if (!hasUppercase || !hasLowercase || !hasNumber) {
    return { 
      isValid: false, 
      message: 'Password must contain at least one uppercase letter, one lowercase letter, and one number' 
    };
  }
  
  return { isValid: true };
}

export function validateUsername(username: string): ValidationResult {
  if (!username) {
    return { isValid: false, message: 'Username is required' };
  }
  
  if (username.length < 3) {
    return { isValid: false, message: 'Username must be at least 3 characters long' };
  }
  
  if (username.length > 30) {
    return { isValid: false, message: 'Username must be less than 30 characters' };
  }
  
  const usernameRegex = /^[a-zA-Z0-9_-]+$/;
  if (!usernameRegex.test(username)) {
    return { 
      isValid: false, 
      message: 'Username can only contain letters, numbers, underscores, and hyphens' 
    };
  }
  
  return { isValid: true };
}

export function validateProjectName(name: string): ValidationResult {
  if (!name) {
    return { isValid: false, message: 'Project name is required' };
  }
  
  if (name.length > APP_CONSTANTS.VALIDATION.MAX_PROJECT_NAME_LENGTH) {
    return { 
      isValid: false, 
      message: `Project name must be less than ${APP_CONSTANTS.VALIDATION.MAX_PROJECT_NAME_LENGTH} characters` 
    };
  }
  
  return { isValid: true };
}

export function validateProjectDescription(description: string): ValidationResult {
  if (description && description.length > APP_CONSTANTS.VALIDATION.MAX_DESCRIPTION_LENGTH) {
    return { 
      isValid: false, 
      message: `Description must be less than ${APP_CONSTANTS.VALIDATION.MAX_DESCRIPTION_LENGTH} characters` 
    };
  }
  
  return { isValid: true };
}

export function validateProjectTags(tags: string[]): ValidationResult {
  if (tags.length > APP_CONSTANTS.VALIDATION.MAX_TAGS_COUNT) {
    return { 
      isValid: false, 
      message: `Maximum ${APP_CONSTANTS.VALIDATION.MAX_TAGS_COUNT} tags allowed` 
    };
  }
  
  for (const tag of tags) {
    if (tag.length > 20) {
      return { isValid: false, message: 'Tag length must be less than 20 characters' };
    }
    
    if (!/^[a-zA-Z0-9_-]+$/.test(tag)) {
      return { 
        isValid: false, 
        message: 'Tags can only contain letters, numbers, underscores, and hyphens' 
      };
    }
  }
  
  return { isValid: true };
}

export function validateRequired(value: any, fieldName: string): ValidationResult {
  if (value === null || value === undefined || value === '') {
    return { isValid: false, message: `${fieldName} is required` };
  }
  
  return { isValid: true };
}
