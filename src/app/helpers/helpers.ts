export const onlyUnique = (value, index, self) => self.indexOf(value) === index;
export const getRandomString = (length = 10): string => {
  let result = '';
  const characters = '!@#$%^&*()_+-=,<.>ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
  const charactersLength = characters.length;
  for (let i = 0; i < length; i++) {
    result += characters.charAt(Math.floor(Math.random() * charactersLength));
  }
  return result;
};
