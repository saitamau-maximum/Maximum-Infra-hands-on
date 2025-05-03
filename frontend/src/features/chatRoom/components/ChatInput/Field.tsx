import { ReactNode } from "react"

import styles from "./Field.module.css"

type FieldProps = {
  children: ReactNode
}

export const Field = ({ children }: FieldProps) => {
  return (
    <div className={styles.field}>
      {children}
    </div>
  )
}