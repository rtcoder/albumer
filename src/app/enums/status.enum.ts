export enum StatusEnum {
  ORDERED = 0,
  OWNED = 1,
}

export const StatusEnumNames = [
  {
    id: StatusEnum.OWNED,
    name: 'Posiadane'
  },
  {
    id: StatusEnum.ORDERED,
    name: 'Zamówione'
  }
];

export const getStatusName = (status: StatusEnum): string => {
  const found = StatusEnumNames.find(el => status === el.id);
  return found ? found.name : '';
};
