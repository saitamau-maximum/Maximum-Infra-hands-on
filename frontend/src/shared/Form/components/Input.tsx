import styles from './Input.module.css';

type InputType = 'text' | 'email' | 'password' | 'number';

type InputProps = {
  type?: InputType;
  id?: string;
  name?: string;
  required?: boolean;
  placeholder?: string;
  minlength?: number;
  maxlength?: number;
}

export const Input = ({
  type,
  id,
  name,
  required,
  placeholder,
  minlength,
  maxlength,
}: InputProps) => {
  return (
    <input
      className={styles.input}
      type={type}
      id={id}
      name={name}
      required={required}
      placeholder={placeholder}
      minLength={minlength}
      maxLength={maxlength}
    />
  );
}