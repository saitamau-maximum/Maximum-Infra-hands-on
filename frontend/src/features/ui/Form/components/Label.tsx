import styles from './Label.module.css'

type LabelProps = {
  label: string
}

export const Label = ({label}: LabelProps) => {
  return (
    <label className={styles.label}>
      {label}
    </label>
  )
}