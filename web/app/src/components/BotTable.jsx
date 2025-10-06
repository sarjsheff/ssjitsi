import LoadingSpinner from './LoadingSpinner';

/**
 * Компонент таблицы ботов
 */
const BotTable = ({
  bots,
  loading,
  error,
  onRefreshScreenshot
}) => {
  const formatTime = (date) => {
    if (!date) return 'Никогда';
    return new Date(date).toLocaleTimeString('ru-RU');
  };

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

  // Сортируем ботов по имени комнаты
  const sortedBots = [...bots].sort((a, b) => {
    const roomA = a.room || '';
    const roomB = b.room || '';
    return roomA.localeCompare(roomB, 'ru');
  });

  return (
    <div className="table-responsive">
      <table className="table table-hover table-sm align-middle">
        <thead className="table-light sticky-top">
          <tr>
            <th style={{ width: '120px' }}>Скриншот</th>
            <th>Комната</th>
            <th>Бот</th>
            <th>Сервер</th>
            <th style={{ width: '100px' }}>Авторизация</th>
            <th style={{ width: '100px' }}>Обновлено</th>
            <th style={{ width: '80px' }}>Статус</th>
            <th style={{ width: '80px' }}>Действия</th>
          </tr>
        </thead>
        <tbody>
          {sortedBots.map(bot => (
            <tr key={bot.id}>
              <td>
                {bot.loadingScreenshot ? (
                  <div className="d-flex justify-content-center align-items-center" style={{ height: '60px' }}>
                    <div className="spinner-border spinner-border-sm" role="status"></div>
                  </div>
                ) : bot.screenshotError ? (
                  <div className="text-center" style={{ height: '60px' }}>
                    <i className="bi bi-exclamation-triangle text-warning"></i>
                  </div>
                ) : bot.screenshot ? (
                  <img
                    src={bot.screenshot}
                    alt={`Скриншот ${bot.id}`}
                    className="img-thumbnail"
                    style={{
                      maxHeight: '60px',
                      width: 'auto',
                      objectFit: 'contain'
                    }}
                  />
                ) : (
                  <div className="text-center text-muted" style={{ height: '60px' }}>
                    <i className="bi bi-image"></i>
                  </div>
                )}
              </td>
              <td>
                <strong>{bot.room}</strong>
              </td>
              <td>{bot.botName}</td>
              <td>
                <span className="text-truncate d-inline-block" style={{ maxWidth: '200px' }} title={bot.server}>
                  {bot.server}
                </span>
              </td>
              <td>
                <span className={`badge ${bot.authMethod === 'JWT' ? 'bg-success' : 'bg-primary'}`}>
                  {bot.authMethod}
                </span>
              </td>
              <td>
                <small className="text-muted">
                  {bot.lastScreenshotUpdate ? formatTime(bot.lastScreenshotUpdate) : 'Ожидание...'}
                </small>
              </td>
              <td>
                <div className="d-flex align-items-center">
                  <div
                    className={`badge ${bot.screenshot ? 'bg-success' : 'bg-secondary'} me-1`}
                    style={{ width: '8px', height: '8px', borderRadius: '50%' }}
                  ></div>
                  <small className="text-nowrap">
                    {bot.screenshot ? 'Активен' : 'Неактивен'}
                  </small>
                </div>
              </td>
              <td>
                <button
                  className="btn btn-sm btn-outline-primary"
                  onClick={() => onRefreshScreenshot(bot.id)}
                  disabled={bot.loadingScreenshot}
                  title="Обновить скриншот"
                >
                  <i className="bi bi-arrow-clockwise"></i>
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default BotTable;
