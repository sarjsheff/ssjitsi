// API клиент для взаимодействия с сервером ботов
const API_BASE_URL = 'http://localhost:8080/api/v1';

/**
 * Получить список ID запущенных ботов
 * @returns {Promise<string[]>}
 */
export const getBots = async () => {
  try {
    const response = await fetch(`${API_BASE_URL}/bots`);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Ошибка при получении списка ботов:', error);
    throw error;
  }
};

/**
 * Получить скриншот бота
 * @param {string} botId - ID бота
 * @returns {Promise<string>} Data URL скриншота
 */
export const getScreenshot = async (botId) => {
  try {
    const response = await fetch(`${API_BASE_URL}/${botId}/screenshot`);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    const blob = await response.blob();
    return URL.createObjectURL(blob);
  } catch (error) {
    console.error(`Ошибка при получении скриншота бота ${botId}:`, error);
    throw error;
  }
};

/**
 * Получить HTML страницы бота
 * @param {string} botId - ID бота
 * @returns {Promise<string>}
 */
export const getBotHtml = async (botId) => {
  try {
    const response = await fetch(`${API_BASE_URL}/${botId}/html`);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error(`Ошибка при получении HTML бота ${botId}:`, error);
    throw error;
  }
};

/**
 * Проверить доступность сервера
 * @returns {Promise<boolean>}
 */
export const checkServerStatus = async () => {
  try {
    const response = await fetch(`${API_BASE_URL}/bots`);
    return response.ok;
  } catch (error) {
    return false;
  }
};