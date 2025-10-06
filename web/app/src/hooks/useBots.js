import { useState, useEffect, useCallback } from 'react';
import { getBots, getScreenshot } from '../services/api';

/**
 * Кастомный хук для управления состоянием ботов
 */
export const useBots = () => {
  const [bots, setBots] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [lastUpdate, setLastUpdate] = useState(null);

  // Функция для загрузки списка ботов
  const loadBots = useCallback(async () => {
    try {
      setError(null);
      const botInfos = await getBots();

      // Создаем или обновляем состояние ботов
      setBots(prevBots => {
        const updatedBots = botInfos.map(botInfo => {
          const existingBot = prevBots.find(bot => bot.id === botInfo.id);
          return existingBot ? {
            ...existingBot,
            room: botInfo.room,
            botName: botInfo.botName,
            server: botInfo.server,
            authMethod: botInfo.authMethod,
            lastUpdate: botInfo.lastUpdate
          } : {
            id: botInfo.id,
            room: botInfo.room,
            botName: botInfo.botName,
            server: botInfo.server,
            authMethod: botInfo.authMethod,
            lastUpdate: botInfo.lastUpdate,
            screenshot: null,
            loadingScreenshot: false,
            screenshotError: null,
            lastScreenshotUpdate: null
          };
        });

        // Удаляем ботов, которых больше нет в списке
        const currentBotIds = botInfos.map(b => b.id);
        return updatedBots.filter(bot => currentBotIds.includes(bot.id));
      });

      setLastUpdate(new Date());
    } catch (err) {
      setError(`Ошибка загрузки ботов: ${err.message}`);
      console.error('Ошибка при загрузке ботов:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  // Функция для загрузки скриншота конкретного бота
  const loadScreenshot = useCallback(async (botId) => {
    try {
      setBots(prevBots => 
        prevBots.map(bot => 
          bot.id === botId 
            ? { ...bot, loadingScreenshot: true, screenshotError: null }
            : bot
        )
      );

      const screenshotUrl = await getScreenshot(botId);
      
      setBots(prevBots => 
        prevBots.map(bot => 
          bot.id === botId 
            ? { 
                ...bot, 
                screenshot: screenshotUrl, 
                loadingScreenshot: false,
                lastScreenshotUpdate: new Date()
              }
            : bot
        )
      );
    } catch (err) {
      setBots(prevBots => 
        prevBots.map(bot => 
          bot.id === botId 
            ? { 
                ...bot, 
                loadingScreenshot: false, 
                screenshotError: `Ошибка загрузки скриншота: ${err.message}` 
              }
            : bot
        )
      );
      console.error(`Ошибка при загрузке скриншота бота ${botId}:`, err);
    }
  }, []);

  // Функция для загрузки скриншотов всех ботов
  const loadAllScreenshots = useCallback(async () => {
    const currentBots = bots;
    for (const bot of currentBots) {
      await loadScreenshot(bot.id);
      // Небольшая задержка между запросами, чтобы не перегружать сервер
      await new Promise(resolve => setTimeout(resolve, 100));
    }
  }, [bots, loadScreenshot]);

  // Функция для обновления всех данных
  const refreshAll = useCallback(async () => {
    setLoading(true);
    await loadBots();
    await loadAllScreenshots();
  }, [loadBots, loadAllScreenshots]);

  // Автоматическое обновление списка ботов
  useEffect(() => {
    loadBots();
    
    const interval = setInterval(loadBots, 10000); // Обновление каждые 10 секунд
    return () => clearInterval(interval);
  }, [loadBots]);

  // Автоматическое обновление скриншотов
  useEffect(() => {
    if (bots.length > 0) {
      const interval = setInterval(loadAllScreenshots, 30000); // Обновление каждые 30 секунд
      return () => clearInterval(interval);
    }
  }, [bots.length, loadAllScreenshots]);

  return {
    bots,
    loading,
    error,
    lastUpdate,
    refreshAll,
    loadScreenshot,
    loadAllScreenshots
  };
};