import LoadingSpinner from './LoadingSpinner';

/**
 * Компонент карточки бота
 */
const BotCard = ({ 
  bot, 
  onRefreshScreenshot 
}) => {
  const formatTime = (date) => {
    if (!date) return 'Никогда';
    return new Date(date).toLocaleTimeString('ru-RU');
  };

  const formatDate = (date) => {
    if (!date) return '';
    return new Date(date).toLocaleDateString('ru-RU');
  };

  return (
    <div className="card h-100 shadow-sm">
      {/* Заголовок карточки */}
      <div className="card-header bg-light">
        <div className="d-flex justify-content-between align-items-center">
          <h6 className="card-title mb-0 text-truncate" title={bot.id}>
            Бот: {bot.id.substring(0, 8)}...
          </h6>
          <button
            className="btn btn-sm btn-outline-primary"
            onClick={() => onRefreshScreenshot(bot.id)}
            disabled={bot.loadingScreenshot}
            title="Обновить скриншот"
          >
            <i className="bi bi-arrow-clockwise"></i>
          </button>
        </div>
      </div>

      <div className="card-body d-flex flex-column">
        {/* Скриншот */}
        <div className="mb-3 text-center">
          {bot.loadingScreenshot ? (
            <div className="d-flex justify-content-center align-items-center" style={{ height: '150px' }}>
              <LoadingSpinner size="sm" text="Загрузка скриншота..." />
            </div>
          ) : bot.screenshotError ? (
            <div className="alert alert-warning text-center py-2" style={{ height: '150px', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
              <div>
                <i className="bi bi-exclamation-triangle fs-4"></i>
                <p className="mb-0 small">{bot.screenshotError}</p>
              </div>
            </div>
          ) : bot.screenshot ? (
            <img
              src={bot.screenshot}
              alt={`Скриншот бота ${bot.id}`}
              className="img-fluid rounded border"
              style={{ 
                maxHeight: '150px', 
                width: 'auto',
                objectFit: 'contain'
              }}
              onError={(e) => {
                e.target.style.display = 'none';
                e.target.nextSibling.style.display = 'block';
              }}
            />
          ) : (
            <div className="text-muted d-flex justify-content-center align-items-center" style={{ height: '150px' }}>
              <div>
                <i className="bi bi-image fs-1"></i>
                <p className="mb-0 small">Скриншот не загружен</p>
              </div>
            </div>
          )}
        </div>

        {/* Информация о боте */}
        <div className="mt-auto">
          <div className="row small text-muted">
            <div className="col-12 mb-1">
              <strong>ID:</strong> 
              <span className="text-truncate d-block" title={bot.id}>
                {bot.id.substring(0, 16)}...
              </span>
            </div>
            
            <div className="col-12 mb-1">
              <strong>Последнее обновление:</strong>
              <div>
                {formatTime(bot.lastScreenshotUpdate)}
                {bot.lastScreenshotUpdate && (
                  <span className="d-block">{formatDate(bot.lastScreenshotUpdate)}</span>
                )}
              </div>
            </div>

            <div className="col-12">
              <div className="d-flex align-items-center">
                <div 
                  className={`badge ${bot.screenshot ? 'bg-success' : 'bg-secondary'} me-2`}
                  style={{ width: '8px', height: '8px', borderRadius: '50%' }}
                ></div>
                <small>
                  {bot.screenshot ? 'Активен' : 'Неактивен'}
                </small>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Футер карточки */}
      <div className="card-footer bg-transparent border-top-0">
        <small className="text-muted">
          {bot.lastScreenshotUpdate 
            ? `Обновлено: ${formatTime(bot.lastScreenshotUpdate)}` 
            : 'Ожидание данных...'
          }
        </small>
      </div>
    </div>
  );
};

export default BotCard;