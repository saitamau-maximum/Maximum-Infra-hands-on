import { forwardRef } from "react";
import styles from './Input.module.css';

type InputType = 'text' | 'email' | 'password' | 'number' | 'file';

type InputProps = {
  type?: InputType;
  id?: string;
  name?: string;
  required?: boolean;
  placeholder?: string;
  minLength?: number;
  maxLength?: number;
} & React.InputHTMLAttributes<HTMLInputElement>;

export const Input = forwardRef<HTMLInputElement, InputProps>(
  (
    {
      type = "text",
      id,
      name,
      required,
      placeholder,
      minLength,
      maxLength,
      ...rest
    },
    ref
  ) => {
    return (
      <input
        ref={ref}
        className={styles.input}
        type={type}
        id={id}
        name={name}
        required={required}
        placeholder={placeholder}
        minLength={minLength}
        maxLength={maxLength}
        {...rest} // register()から渡されるonChangeやonBlurなどを渡す
      />
    );
  }
);

Input.displayName = "Input";
