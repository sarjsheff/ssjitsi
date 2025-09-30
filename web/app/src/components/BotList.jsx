import BotCard from './BotCard';
import LoadingSpinner from './LoadingSpinner';

/**
 * Компонент списка ботов
 */
const BotList = ({ 
  bots, 
  loading, 
  error, 
  onRefreshScreenshot 
}) => {
  if (loading && bots.length === 0) {
    return (
      <div className="d-flex justify-content-center align-items-center py-5">
        <LoadingSpinner size="lg" text="Загрузка ботов..." />
      </div>
    );
  }

  if (error) {
    return (
      <div className="alert alert-danger text-center">
        <i className="bi bi-exclamation-triangle-fill me-2"></i>
        {error}
      </div>
    );
  }

  if (bots.length === 0) {
    return (
      <div className="text-center py-5">
        <i className="bi bi-robot fs-1 text-muted"></i>
        <h5 className="text-muted mt-3">Боты не найдены</h5>
        <p className="text-muted">Запустите сервер и настройте ботов в конфигурации</p>
      </div>
    );
  }

  return (
    <div className="row g-3">
      {bots.map(bot => (
        <div key={bot.id} className="col-12 col-sm-6 col-md-4 col-lg-3">
          <BotCard 
            bot={bot} 
            onRefreshScreenshot={onRefreshScreenshot}
          />
        </div>
      ))}
    </div>
  );
};

export default BotList;