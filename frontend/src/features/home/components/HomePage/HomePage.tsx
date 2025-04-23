import styles from './HomePage.module.css';

export const HomePage = () => {
  return (
    <div className={styles.container}>
      <h1 className={styles.hero}>Chat-INFRA</h1>
      <p>ようこそ、ユーザーさん！</p>
    </div>
  );
};
