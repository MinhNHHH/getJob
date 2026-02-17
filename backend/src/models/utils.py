############################### Helpers functions ###############################

def insert(db, model, obj):
  """
  Insert a record then commit.

  >> create(db, TwitterAuthor, {"username": "qnkhuat", "description": "Who Am I, where I am?"})

  Parameters
  ----------

    db: Session
      Database session

    model: Model
      The SqlAlchemy Model

    obj_in: dict
      A dictionary with values for the record

  Returns
  -------
    Returns the created db_obj
  """
  db_obj = model(**obj)
  db.add(db_obj)
  db.commit()
  db.refresh(db_obj)
  return db_obj

def update(db, db_obj, obj_in):
  """
  Update a sqlalchemy object and commit it.

  >> author = db.query(TwitterAuthor).filter(TwitterAuthor.username == "qnkhuat").first()
  >> updated_author = update(db, author, {"description": "I code, therefore I am."})

  Parameters
  ----------

    db: Session
      Database session

    db_obj: Model
      The SqlAlchemy database object

    obj_in: dict
      A dictionary with values to update

  Returns
  -------
    Returns the updated db_obj
  """
  for field in db_obj.__mapper__.attrs.keys():
    if field in obj_in:
      setattr(db_obj, field, obj_in[field])
  db.add(db_obj)
  db.commit()
  db.refresh(db_obj)
  return db_obj

def delete(db, db_obj):
  """
  Delete a database object.
  """
  db.delete(db_obj)
  db.commit()


def row2dict(obj):
  """
  Turns an SQLAchemy object row to a dict.
  row =  db.query(models.TwitterAuthor.id, models.TwitterAuthor.username).one()

  row2dict(row)
  => {
  "id":       1,
  "username": "qnkhuat"
  }
  """
  return {key: obj[key] for key in obj.keys()}