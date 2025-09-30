/**
 * Компонент индикатора загрузки
 */
const LoadingSpinner = ({ size = 'sm', text = 'Загрузка...' }) => {
  const sizeClass = {
    sm: 'spinner-border-sm',
    md: '',
    lg: 'spinner-border-lg'
  }[size];

  return (
    <div className="d-flex align-items-center justify-content-center">
      <div className={`spinner-border ${sizeClass} text-primary`} role="status">
        <span className="visually-hidden">{text}</span>
      </div>
      {text && <span className="ms-2">{text}</span>}
    </div>
  );
};

export default LoadingSpinner;