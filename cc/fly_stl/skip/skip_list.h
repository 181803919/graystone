#ifndef _SKIP_LIST_H_
#define _SKIP_LIST_H_

constexpr int NODE_UP_RADE = 5000;
enum ENUM_SKIP_LIST_NODE_TYPE
{
	SKIP_LIST_NODE_HEAD = 1,
	SKIP_LIST_NODE_NORMAL = 2,
	SKIP_LIST_NODE_TAIL = 3,
};

/*--------------------------------------------Skip_List_Node-------------------------------------*/
template<class T, class Compare_Fun, class Equal_Fun>
class Skip_List_Node
{
public:
	Skip_List_Node();
	Skip_List_Node(T &data);
	~Skip_List_Node();

public:
	void Reinit();
	void Reinit(T& data);

public:
	T data_;
	Skip_List_Node* up_;
	Skip_List_Node* down_;
	Skip_List_Node* left_;
	Skip_List_Node* right_;
	ENUM_SKIP_LIST_NODE_TYPE node_type_;
};

template<class T, class Compare_Fun, class Equal_Fun>
Skip_List_Node<T, Compare_Fun, Equal_Fun>::Skip_List_Node()
{
	memset(&data_, 0, sizeof(data_));
	node_type_ = SKIP_LIST_NODE_NORMAL;
	up_ = NULL;
	down_ = NULL;
	left_ = NULL;
	right_ = NULL;
}

template<class T, class Compare_Fun, class Equal_Fun>
Skip_List_Node<T, Compare_Fun, Equal_Fun>::Skip_List_Node(T& data)
{
	node_type_ = SKIP_LIST_NODE_NORMAL;
	data_ = data;
	up_ = NULL;
	down_ = NULL;
	left_ = NULL;
	right_ = NULL;
}

template<class T, class Compare_Fun, class Equal_Fun>
Skip_List_Node<T, Compare_Fun, Equal_Fun>::~Skip_List_Node()
{

}

template<class T, class Compare_Fun, class Equal_Fun>
void Skip_List_Node<T, Compare_Fun, Equal_Fun>::Reinit()
{
	node_type_ = SKIP_LIST_NODE_NORMAL;
	memset(&data_, 0, sizeof(data_));
	up_ = NULL;
	down_ = NULL;
	left_ = NULL;
	right_ = NULL;
}

template<class T, class Compare_Fun, class Equal_Fun>
void Skip_List_Node<T, Compare_Fun, Equal_Fun>::Reinit(T &data)
{
	node_type_ = SKIP_LIST_NODE_NORMAL;
	data_ = data;
	up_ = NULL;
	down_ = NULL;
	left_ = NULL;
	right_ = NULL;
}


/*-----------------------------------------------Skip_List---------------------------------------*/
typedef void (*SkipListForeachFun)(void* data);

template<class T, class Compare_Fun, class Equal_Fun>
class Skip_List
{
public:
	Skip_List();
	~Skip_List();

public:
	typedef Skip_List_Node<T, Compare_Fun, Equal_Fun> Node;

public:
	Node* GetMaxLevelHead();

	Node* SeachNode(T& data);
	Node* FindPreNodeInRow(Node* row_head, T& data);
	Node* FindPreNode(T& data);

	Node* InsertNodeInRow(Node* row_head, T& data, Node* down_node = NULL);
	Node* InsertNode(T& data);

	Node* DeleteNodeInRow(Node* node);
	void DeleteNode(T& data);

	void IncNum();
	void DecNum();
	bool IfNeedUp(Node* new_node = NULL);

	void  AddRow();
	void  DelRow();

	void  Foreach(SkipListForeachFun fun);

	Node* AllocNode(ENUM_SKIP_LIST_NODE_TYPE type);
	Node* AllocNode(T &data);

	void  BackNode(Node* node);

	void  SetDebug(bool if_debug);
	void  SetPrintFun(SkipListForeachFun print_fun);

	int   Reserve(int resize_num);
public:
	Node* head_;
	Node* tail_;
	Node* top_head_;
	Node* top_tail_;

	Compare_Fun compare_fun_;
	Equal_Fun	equal_fun_;
	int level_;
	int node_num_;

	bool if_debug_;
	SkipListForeachFun print_fun_;

	int pre_molloc_node_num_;
	std::queue<Node*> mem_list_;

	int insert_deep_;
};

template<class T, class Compare_Fun, class Equal_Fun>
Skip_List<T, Compare_Fun, Equal_Fun>::Skip_List()
{
	pre_molloc_node_num_ = 0;
	head_ = AllocNode(SKIP_LIST_NODE_HEAD);
	tail_ = AllocNode(SKIP_LIST_NODE_TAIL);

	head_->right_ = tail_;
	tail_->left_ = head_;
	top_head_ = head_;
	top_tail_ = tail_;

	level_ = 1;
	node_num_ = 0;
	if_debug_ = false;
	print_fun_ = NULL;
}

template<class T, class Compare_Fun, class Equal_Fun>
Skip_List<T, Compare_Fun, Equal_Fun>::~Skip_List()
{
	//if (pre_molloc_node_num_ == 0)
	{
		Node* ver_node = top_head_;
		while (NULL != ver_node)
		{
			Node* hor_node = ver_node;
			ver_node = ver_node->down_;
			while (NULL != hor_node)
			{
				Node* tmp = hor_node;
				hor_node = hor_node->right_;
				delete tmp;
			}
		}
	}
	//else
	//{
	//	Node* del_node = NULL;
	//	while (true)
	//	{
	//		del_node = mem_list_.front();
	//		mem_list_.pop();
	//		if (del_node == NULL)
	//		{
	//			break;
	//		}
	//		delete del_node;
	//		del_node = NULL;
	//	}
	//}
}

template<class T, class Compare_Fun, class Equal_Fun>
Skip_List_Node<T, Compare_Fun, Equal_Fun>*
Skip_List<T, Compare_Fun, Equal_Fun>::GetMaxLevelHead()
{
	return top_head_;
}

template<class T, class Compare_Fun, class Equal_Fun>
Skip_List_Node<T, Compare_Fun, Equal_Fun>*
Skip_List<T, Compare_Fun, Equal_Fun>::SeachNode(T& data)
{
	Node* pre_node = FindPreNode(data);
	if (pre_node == NULL)
	{
		return NULL;
	}

	if (pre_node->right_ == NULL)
	{
		return NULL;
	}

	if (pre_node->right_->node_type_ != SKIP_LIST_NODE_NORMAL)
	{
		return NULL;
	}

	if (!equal_fun_(pre_node->right_->data_, data))
	{
		return NULL;
	}
	return pre_node->right_;
}

template<class T, class Compare_Fun, class Equal_Fun>
Skip_List_Node<T, Compare_Fun, Equal_Fun>*
Skip_List<T, Compare_Fun, Equal_Fun>::FindPreNodeInRow(Node* row_head, T& data)
{
	Node* hor_node = row_head;
	while (true)
	{
		hor_node = hor_node->right_;
		if (if_debug_)
		{
			(*print_fun_)(hor_node);
		}

		if (SKIP_LIST_NODE_NORMAL != hor_node->node_type_)
		{
			break;
		}

		if (compare_fun_(data, hor_node->data_))
		{
			break;
		}

		if (equal_fun_(data, hor_node->data_))
		{
			break;
		}
	}

	return hor_node->left_;
}

template<class T, class Compare_Fun, class Equal_Fun>
Skip_List_Node<T, Compare_Fun, Equal_Fun>*
Skip_List<T, Compare_Fun, Equal_Fun>::FindPreNode(T& data)
{
	int cur_level = 1;

	Node* ver_node = top_head_;
	while(true)
	{
		if (if_debug_)
		{
			printf("\nInLevel:%d\n", cur_level);
		}

		ver_node = FindPreNodeInRow(ver_node, data);
		if (ver_node->down_ == NULL)
		{
			break;
		}

		ver_node = ver_node->down_;
		++cur_level;
	}
	
	return ver_node;
}

template<class T, class Compare_Fun, class Equal_Fun>
Skip_List_Node<T, Compare_Fun, Equal_Fun>*
Skip_List<T, Compare_Fun, Equal_Fun>::InsertNodeInRow(Node* row_head, T& data, Node* down_node)
{
	Node* new_node = AllocNode(data);
	if (new_node == NULL)
	{
		return NULL;
	}

	IncNum();
	Node* pre_node = FindPreNodeInRow(row_head, data);
	new_node->left_ = pre_node;
	new_node->right_ = pre_node->right_;
	new_node->right_->left_ = new_node;
	new_node->left_->right_ = new_node;

	if (NULL != down_node)
	{
		new_node->down_ = down_node;
		down_node->up_ = new_node;
	}

	return new_node;
}

template<class T, class Compare_Fun, class Equal_Fun>
Skip_List_Node<T, Compare_Fun, Equal_Fun>*
Skip_List<T, Compare_Fun, Equal_Fun>::InsertNode(T& data)
{
	Node* pre_node = FindPreNode(data);
	if (pre_node->right_->node_type_ == SKIP_LIST_NODE_NORMAL
		&& equal_fun_(pre_node->right_->data_, data))
	{
		return pre_node->right_;
	}

	int cur_lvl = 0;

	insert_deep_ = 0;
	bool have_add = false;
	Node* ver_node = pre_node;
	Node* new_node = NULL;
	while (true)
	{
		++cur_lvl;

		new_node = InsertNodeInRow(ver_node, data, new_node);
		if (new_node == NULL)
		{
			return NULL;
		}

		if (IfNeedUp() && !have_add)
		{
			if (cur_lvl == level_)
			{
				AddRow();
				have_add = true;
			}
			
			Node* left = new_node->left_;
			while (true)
			{
				if (left->up_ != NULL)
				{
					ver_node = left->up_;
					break;
				}
				left = left->left_;
			}
		}
		else
		{
			break;
		}
	}
	return new_node;
}

template<class T, class Compare_Fun, class Equal_Fun>
Skip_List_Node<T, Compare_Fun, Equal_Fun> *
Skip_List<T, Compare_Fun, Equal_Fun>::DeleteNodeInRow(Node* node)
{
	Node* up = node->up_;
	Node* left = node->left_;
	Node* right = node->right_;

	left->right_ = right;
	right->left_ = left;
	BackNode(node);

	while(top_head_->right_->node_type_ == SKIP_LIST_NODE_TAIL
		&& level_ > 1)
	{
		DelRow();
	}

	return up;
}

template<class T, class Compare_Fun, class Equal_Fun>
void Skip_List<T, Compare_Fun, Equal_Fun>::DeleteNode(T& data)
{
	Node* node = SeachNode(data);
	if (node == NULL)
	{
		return;
	}

	while (true)
	{
		node = DeleteNodeInRow(node);
		if (NULL == node)
		{
			break;
		}
	}
}

template<class T, class Compare_Fun, class Equal_Fun>
void Skip_List<T, Compare_Fun, Equal_Fun>::IncNum()
{
	++node_num_;
}

template<class T, class Compare_Fun, class Equal_Fun>
void Skip_List<T, Compare_Fun, Equal_Fun>::DecNum()
{
	--node_num_;
}

template<class T, class Compare_Fun, class Equal_Fun>
bool Skip_List<T, Compare_Fun, Equal_Fun>::IfNeedUp(Node* new_node)
{
	if (new_node == NULL)
	{
		int value = rand() % 10000;
		return value < NODE_UP_RADE;
	}
	else
	{
		void* p = reinterpret_cast<void*>(new_node);
		long long value = (long long)(p);
		return value % 2 == 0;
	}
}

template<class T, class Compare_Fun, class Equal_Fun>
void Skip_List<T, Compare_Fun, Equal_Fun>::AddRow()
{
	Node* new_head = AllocNode(SKIP_LIST_NODE_HEAD);
	Node* new_tail = AllocNode(SKIP_LIST_NODE_TAIL);
	new_head->right_ = new_tail;
	new_tail->left_ = new_head;

	new_head->down_ = top_head_;
	new_tail->down_ = top_tail_;

	top_head_->up_ = new_head;
	top_tail_->up_ = new_tail;

	top_head_ = new_head;
	top_tail_ = new_tail;

	++level_;
}

template<class T, class Compare_Fun, class Equal_Fun>
void Skip_List<T, Compare_Fun, Equal_Fun>::DelRow()
{
	Node* top_head = top_head_;
	Node* top_tail = top_tail_;

	top_head_ = top_head->down_;
	top_tail_ = top_tail->down_;

	top_head_->up_ = NULL;
	top_tail_->up_ = NULL;

	BackNode(top_head);
	BackNode(top_tail);

	--level_;
}

template<class T, class Compare_Fun, class Equal_Fun>
Skip_List_Node<T, Compare_Fun, Equal_Fun>* 
Skip_List<T, Compare_Fun, Equal_Fun>::AllocNode(ENUM_SKIP_LIST_NODE_TYPE type)
{
	Node* new_node = NULL;
	if (pre_molloc_node_num_ == 0)
	{
		new_node = new Node();
	}
	else
	{
		new_node = mem_list_.front();
		mem_list_.pop();
		new_node->Reinit();
	}

	new_node->node_type_ = type;
	return new_node;
}

template<class T, class Compare_Fun, class Equal_Fun>
Skip_List_Node<T, Compare_Fun, Equal_Fun>*
Skip_List<T, Compare_Fun, Equal_Fun>::AllocNode(T &data)
{
	Node* new_node = NULL;
	if (pre_molloc_node_num_ == 0)
	{
		new_node = new Node(data);
	}
	else
	{
		new_node = mem_list_.front();
		mem_list_.pop();
		new_node->Reinit(data);
	}

	new_node->node_type_ = SKIP_LIST_NODE_NORMAL;
	return new_node;
}

template<class T, class Compare_Fun, class Equal_Fun>
void Skip_List<T, Compare_Fun, Equal_Fun>::BackNode(Node* node)
{
	if (pre_molloc_node_num_ == 0)
	{
		delete node;
	}
	else
	{
		mem_list_.push(node);
	}
}

template<class T, class Compare_Fun, class Equal_Fun>
void Skip_List<T, Compare_Fun, Equal_Fun>::Foreach(SkipListForeachFun fun)
{
	Node* ver_node = top_head_;
	while (NULL != ver_node)
	{
		Node* hor_node = ver_node;
		ver_node = ver_node->down_;
		while (NULL != hor_node)
		{
			Node* tmp = hor_node;
			hor_node = hor_node->right_;
			(*fun)(tmp);
		}
	}
}

template<class T, class Compare_Fun, class Equal_Fun>
void Skip_List<T, Compare_Fun, Equal_Fun>::SetDebug(bool if_debug)
{
	if_debug_ = if_debug;
}

template<class T, class Compare_Fun, class Equal_Fun>
void Skip_List<T, Compare_Fun, Equal_Fun>::SetPrintFun(SkipListForeachFun print_fun)
{
	print_fun_ = print_fun;
}

template<class T, class Compare_Fun, class Equal_Fun>
int Skip_List<T, Compare_Fun, Equal_Fun>::Reserve(int resize_num)
{
	pre_molloc_node_num_ = resize_num;
	int real_use_node = resize_num * 3;
	for (int i = 0; i < real_use_node; ++i)
	{
		Node* new_node = new Node();
		if (new_node == NULL)
		{
			return -1;
		}

		mem_list_.push(new_node);
	}
	return 0;
}

template<class T>
class Compare_Fun
{
public:
	bool operator()(T a, T b) const
	{
		return a < b;
	}
};

template<class T>
class Equal_Fun
{
public:
	bool operator()(T a, T b) const
	{
		return a == b;
	}
};

void PrintSkipNodeInt(void* data)
{
	typedef Skip_List_Node<int, Compare_Fun<int>, Equal_Fun<int>> Node_Type;
	Node_Type* node = reinterpret_cast<Node_Type*>(data);
	if (node->node_type_ == SKIP_LIST_NODE_HEAD)
	{
		printf("r:");
	}
	else if (node->node_type_ == SKIP_LIST_NODE_NORMAL)
	{
		printf(" %d", node->data_);
	}
	else
	{
		printf("\r\n");
	}
}

#endif // !_SKIP_LIST_H_
